package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"pnas/bt"
	"pnas/category"
	"pnas/crawler"
	"pnas/log"
	"pnas/phttp"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/service/chat"
	"pnas/service/session"
	"pnas/setting"
	"pnas/user"
	"pnas/user/task"
	"pnas/utils"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jinzhu/copier"
)

type CoreService struct {
	prpc.UnimplementedUserServiceServer

	sessions session.ISessions

	um          user.UserManger
	shutDownCtx context.Context
	closeFunc   context.CancelFunc

	wg sync.WaitGroup

	grpcSer   *grpc.Server
	httpSer   *http.Server
	videoSer  *VideoService
	posterSer *PosterService
	rooms     chat.IRooms

	shares IItemShares
}

func (ser *CoreService) Init() {
	ss := &session.Sessions{}
	ser.sessions = ss
	ss.Init()

	ser.um.Init()

	var rooms chat.Rooms
	rooms.Init()
	ser.rooms = &rooms

	sm := &ShareManager{}
	sm.Init()
	ser.shares = sm

	if setting.GS().Server.EnableCrawler {
		go crawler.Go36dmBackgroup(&ser.um, -1)
		go crawler.GoAcgBackgroup(&crawler.GoAcgBackgroupParams{
			MagnetShares: &ser.um,
			MaxDepth:     -1,
			ProxyUrl:     setting.GS().Server.CrawlerProxy,
		})
	}
}

func (ser *CoreService) Serve() {
	creds, err := credentials.NewServerTLSFromFile(setting.GS().Server.CrtFile, setting.GS().Server.KeyFile)
	if err != nil {
		log.Panic(err)
	}

	ser.grpcSer = grpc.NewServer(grpc.Creds(creds))
	prpc.RegisterUserServiceServer(ser.grpcSer, ser)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d",
		setting.GS().Server.BoundIp,
		setting.GS().Server.Port))
	if err != nil {
		log.Panic(err)
	}
	go func() {
		if err := ser.grpcSer.Serve(lis); err != nil {
			log.Panic(err)
		}
	}()

	grpcwebServer := grpcweb.WrapServer(ser.grpcSer)
	router := mux.NewRouter()
	router.Use(phttp.CorsMiddleware)
	router.NotFoundHandler = http.HandlerFunc(phttp.NotFound)

	router.Handle("/prpc.UserService/{method}", grpcwebServer)

	ser.videoSer = newVideoService(&NewVideoServiceParams{
		UserData:    &ser.um,
		Shares:      ser.shares,
		Sessions:    ser.sessions,
		Router:      router.PathPrefix("/video").Subrouter(),
		RecvDanmaku: ser.recvDanmaku,
	})
	ser.posterSer = newPosterService(&NewPosterServiceParams{
		CategoryData: ser.um.CategoryService(),
		Shares:       ser.shares,
		Sessions:     ser.sessions,
		Router:       router.PathPrefix("/poster").Subrouter(),
	})

	ser.httpSer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", setting.GS().Server.BoundIp, setting.GS().Server.WebPort),
		Handler: router,
	}

	ser.shutDownCtx, ser.closeFunc = context.WithCancel(context.Background())

	if err := ser.httpSer.ListenAndServeTLS(setting.GS().Server.CrtFile, setting.GS().Server.KeyFile); err != nil {
		log.Panic(err)
	}
}

func (ser *CoreService) recvDanmaku(vid ptype.VideoID, danmaku string) {
	roomKey := utils.GetItemRoomKey(&prpc.Room{
		Type: prpc.Room_Danmaku,
		Id:   int64(vid),
	})

	r := ser.rooms.GetRoom(roomKey)
	if r == nil {
		ser.rooms.CreateRoom(&chat.CreateRoomParams{
			RoomKey:       roomKey,
			ImmediatePush: false,
			Interval:      time.Second * 3,
		})
		r = ser.rooms.GetRoom(roomKey)
		if r == nil {
			return
		}
	}
	ser.rooms.Broadcast(roomKey, &chat.ChatMessage{
		UserId:   user.AdminId,
		SentTime: time.Now(),
		Msg:      danmaku,
	})
}

func (ser *CoreService) Close() {
	ser.closeFunc()
	ser.httpSer.Shutdown(context.Background())
	ser.grpcSer.GracefulStop()
	ser.wg.Wait()
}

func (ser *CoreService) GetUserManager() *user.UserManger {
	return &ser.um
}

func (ser *CoreService) GetSession(r *http.Request) *session.Session {
	s, _ := ser.sessions.GetSession(r)
	return s
}

func (ser *CoreService) getSession(ctx context.Context) *session.Session {
	s, _ := ser.sessions.GetSession2(ctx)
	return s
}

func (ser *CoreService) Register(
	ctx context.Context,
	registerInfo *prpc.RegisterInfo) (*prpc.RegisterRet, error) {

	if registerInfo == nil || registerInfo.GetUserInfo() == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}

	userInfo := registerInfo.GetUserInfo()
	if len(userInfo.GetEmail()) == 0 ||
		len(userInfo.GetName()) == 0 ||
		len(userInfo.GetPasswd()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "")
	}

	params := &user.NewUserParams{
		Email:           userInfo.GetEmail(),
		Name:            userInfo.GetName(),
		Passwd:          userInfo.GetPasswd(),
		Auth:            utils.NewBitSet(user.AuthMax, user.AuthAdmin),
		CategoryService: ser.um.CategoryService(),
	}
	log.Debugf("[user] %s %s register", params.Name, params.Email)
	err := user.NewUser(params)
	if err != nil {
		log.Errorf("[user] %s failed to register, err: %+v", userInfo.GetName(), err)
		return nil, status.Error(codes.Unknown, "")
	}
	return &prpc.RegisterRet{}, nil
}

func (ser *CoreService) IsUsedEmail(
	ctx context.Context,
	emailInfo *prpc.EmailInfo) (ret *prpc.IsUsedEmailRet, err error) {

	if user.UsedEmail(emailInfo.GetEmail()) {
		return nil, status.Error(codes.AlreadyExists, "existed email")
	}
	return &prpc.IsUsedEmailRet{}, nil
}

func (ser *CoreService) Login(ctx context.Context, loginInfo *prpc.LoginInfo) (*prpc.LoginRet, error) {
	user, err := ser.um.Login(loginInfo.GetEmail(), loginInfo.GetPasswd())
	if err != nil {
		log.Warnf("[user] email %s failed to load user err: %v", loginInfo.GetEmail(), err)
		return nil, status.Error(codes.NotFound, "")
	}
	userInfo := user.GetUserInfo()

	s := ser.getSession(ctx)
	var expiresAt time.Time
	if loginInfo.RememberMe {
		expiresAt = time.Now().Add(time.Hour * 24 * 7)
	}
	if s == nil {
		s = ser.sessions.NewSession(&session.NewSessionParams{
			OldId:     -1,
			ExpiresAt: expiresAt,
			UserId:    userInfo.Id,
		})
	} else {
		s = ser.sessions.NewSession(&session.NewSessionParams{
			OldId:     s.Id,
			ExpiresAt: expiresAt,
			UserId:    userInfo.Id,
		})
	}

	grpc.SendHeader(ctx, metadata.Pairs("Set-Cookie",
		session.GenSessionTokenCookie(s), "Set-Cookie", session.GenSessionIdCookie(s)))

	return &prpc.LoginRet{
		Token: s.Token,
		UserInfo: &prpc.UserInfo{
			Id:              int64(userInfo.Id),
			Name:            userInfo.Name,
			Email:           userInfo.Email,
			HomeDirectoryId: int64(userInfo.HomeDirectoryId),
			MagnetRootId:    int64(ser.um.GetMagnetRootId()),
		},
		RememberMe: loginInfo.RememberMe,
	}, nil
}

func (ser *CoreService) FastLogin(
	ctx context.Context,
	loginInfo *prpc.LoginInfo) (*prpc.LoginRet, error) {

	oldSession := ser.getSession(ctx)
	if oldSession == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}

	user, err := ser.um.LoadUser(oldSession.UserId)
	if err != nil {
		log.Warnf("[user] %d failed to fast login, err: %v", user.GetUserInfo().Id, err)
		return nil, status.Error(codes.NotFound, "")
	}
	userInfo := user.GetUserInfo()
	var expiresAt time.Time
	if loginInfo.RememberMe {
		expiresAt = time.Now().Add(time.Hour * 24 * 7)
	}

	s := ser.sessions.NewSession(&session.NewSessionParams{
		OldId:     oldSession.Id,
		ExpiresAt: expiresAt,
		UserId:    userInfo.Id,
	})

	grpc.SendHeader(ctx, metadata.Pairs("Set-Cookie",
		session.GenSessionTokenCookie(s), "Set-Cookie", session.GenSessionIdCookie(s)))

	return &prpc.LoginRet{
		Token: s.Token,
		UserInfo: &prpc.UserInfo{
			Id:              int64(userInfo.Id),
			Name:            userInfo.Name,
			Email:           userInfo.Email,
			HomeDirectoryId: int64(userInfo.HomeDirectoryId),
			MagnetRootId:    int64(ser.um.GetMagnetRootId()),
		},
		RememberMe: loginInfo.RememberMe,
	}, nil
}

func (ser *CoreService) IsLogined(context.Context, *prpc.LoginInfo) (*prpc.LoginRet, error) {
	return &prpc.LoginRet{}, nil
}

func (ser *CoreService) ChangePassword(ctx context.Context, req *prpc.ChangePasswordReq) (*prpc.ChangePasswordRsp, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}

	b := ser.um.ChangePassword(&user.ChangePasswordParams{
		UserId:      ses.UserId,
		Email:       req.Email,
		OldPassword: req.OldPasswd,
		NewPassword: req.NewPasswd,
	})
	if !b {
		return nil, status.Error(codes.InvalidArgument, "更改失败")
	}
	return &prpc.ChangePasswordRsp{}, nil
}

func (ser *CoreService) Download(
	ctx context.Context,
	req *prpc.DownloadRequest) (*prpc.DownloadRespone, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	params := &bt.DownloadParams{
		UserId: ses.UserId,
		Req:    req,
	}
	return ser.um.Download(params)
}

func (ser *CoreService) RemoveTorrent(ctx context.Context, req *prpc.RemoveTorrentReq) (*prpc.RemoveTorrentRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	params := &bt.RemoveTorrentParams{
		UserId: ses.UserId,
		Req:    req,
	}
	return ser.um.RemoveTorrent(params)
}

func (ser *CoreService) GetMagnetUri(ctx context.Context, req *prpc.GetMagnetUriReq) (*prpc.GetMagnetUriRsp, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	params := &bt.GetMagnetUriParams{
		Req: req,
	}
	return ser.um.GetMagnetUri(params)
}

func (ser *CoreService) GetTorrents(ctx context.Context, req *prpc.GetTorrentsReq) (*prpc.GetTorrentsRsp, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "req == nil")
	}
	ts := ser.um.GetTorrents(ses.UserId)
	rts := []*prpc.TorrentInfo{}
	for _, t := range ts {
		var rt prpc.TorrentInfo
		copier.Copy(&rt, t.GetBaseInfo())
		fs := t.GetFiles()
		copier.Copy(&rt.Files, fs)
		rt.SavePath = ""
		rts = append(rts, &rt)
	}
	rsp := &prpc.GetTorrentsRsp{
		TorrentInfo: rts,
	}
	return rsp, nil
}

func (ser *CoreService) GetPeerInfo(ctx context.Context, req *prpc.GetPeerInfoReq) (*prpc.GetPeerInfoRsp, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "req == nil")
	}
	return ser.um.GetPeerInfo(&bt.GetPeerInfoParams{
		UserId: ses.UserId,
		Req:    req,
	})
}

func (ser *CoreService) OnBtStatus(statusReq *prpc.BtStatusRequest, stream prpc.UserService_OnBtStatusServer) error {
	ses := ser.getSession(stream.Context())
	if ses == nil {
		return status.Error(codes.PermissionDenied, "")
	}

	done, c := context.WithCancel(context.Background())
	onError := func() {
		c()
	}

	ser.um.SetSessionCallback(ses.UserId, ses.Id, func(err error, s *prpc.TorrentStatus) {
		if err != nil {
			return
		}
		ret := &prpc.BtStatusRespone{
			StatusArray: []*prpc.TorrentStatus{s},
		}
		err = stream.Send(ret)
		if err != nil {
			onError()
		}
	})

	select {
	case <-stream.Context().Done():
	case <-done.Done():
	}

	ser.um.SetSessionCallback(ses.UserId, ses.Id, nil)

	return nil
}

func (ser *CoreService) QueryBtVideos(ctx context.Context, req *prpc.QueryBtVideosReq) (*prpc.QueryBtVideosRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	if req.InfoHash == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	ms, err := ser.um.QueryBtVideoMetadata(ses.UserId, bt.TranInfoHash(req.InfoHash))
	if err != nil {
		log.Warnf("[user] user %d infoHash %s, query bt video meta err: %s", ses.UserId, hex.EncodeToString(req.InfoHash.Hash), err.Error())
		return nil, err
	}
	res := &prpc.QueryBtVideosRes{}
	for i, m := range ms {
		var rvm prpc.VideoMetadata
		copier.Copy(&rvm, m)
		var meta prpc.BtFileMetadata
		meta.FileIndex = int32(i)
		meta.Meta = &rvm
		res.Data = append(res.Data, &meta)
	}
	return res, nil
}

func (ser *CoreService) NewCategoryItem(ctx context.Context, req *prpc.NewCategoryItemReq) (*prpc.NewCategoryItemRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	params := &category.NewCategoryParams{
		ParentId:     ptype.CategoryID(req.ParentId),
		Creator:      ses.UserId,
		TypeId:       prpc.CategoryItem_Directory,
		Name:         req.Name,
		ResourcePath: req.ResourcePath,
		Introduce:    req.Introduce,
		Auth:         utils.NewBitSet(category.AuthMax),
	}
	err := ser.um.NewCategoryItem(ses.UserId, params)
	if err != nil {
		log.Warnf("[user] %d new category err: %v", ses.UserId, err)
		return nil, status.Error(codes.Unknown, "")
	}
	return &prpc.NewCategoryItemRes{}, nil
}

func (ser *CoreService) DelCategoryItem(ctx context.Context, req *prpc.DelCategoryItemReq) (*prpc.DelCategoryItemRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	err := ser.um.DelCategoryItem(ses.UserId, ptype.CategoryID(req.GetItemId()))
	if err != nil {
		log.Warnf("[user] %d delete category %d err: %v", ses.UserId, req.GetItemId(), err)
		return nil, status.Error(codes.Unknown, "")
	}
	return &prpc.DelCategoryItemRes{}, nil
}

func (ser *CoreService) RenameItem(ctx context.Context, req *prpc.RenameItemReq) (*prpc.RenameItemRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	item, err := ser.um.CategoryService().GetItem(ses.UserId, ptype.CategoryID(req.ItemId))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	if item.GetName() == req.NewName {
		return nil, status.Error(codes.InvalidArgument, "same name")
	}
	err = item.Rename(req.NewName)
	if item.GetId() != category.RootId {
		ser.um.CategoryService().RefreshItem(item.GetParentId())
	}
	if err != nil {
		return nil, err
	}
	return &prpc.RenameItemRes{}, nil
}

func (ser *CoreService) QuerySubItems(ctx context.Context, req *prpc.QuerySubItemsReq) (*prpc.QuerySubItemsRes, error) {
	var userId ptype.UserID
	if len(req.ShareId) > 0 {
		si, err := ser.shares.GetShareItemInfo(req.ShareId)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "not found item")
		}
		if !ser.um.CategoryService().IsRelationOf(ptype.CategoryID(req.ParentId), si.ShareItemInfo.ItemId) {
			return nil, status.Error(codes.PermissionDenied, "not found item")
		}
		userId = si.UserId
	} else {
		ses := ser.getSession(ctx)
		if ses == nil {
			return nil, status.Error(codes.PermissionDenied, "not found session")
		}
		userId = ses.UserId
	}

	parentItem, err := ser.um.CategoryService().GetItem(userId, ptype.CategoryID(req.ParentId))
	if err != nil {
		log.Warnf("[user] %d query category %d err: %v", userId, req.ParentId, err)
		return nil, status.Error(codes.PermissionDenied, "not found")
	}
	items, err := ser.um.CategoryService().GetItemsByParent(
		&category.GetItemsByParentParams{
			Querier:  userId,
			ParentId: ptype.CategoryID(req.ParentId),
			PageNum:  req.PageNum,
			Rows:     req.Rows,
			Desc:     req.Desc,
			Types:    category.NewCategoryTypes(req.Types...),
		})
	if err != nil {
		return nil, err
	}

	var resParentItem prpc.CategoryItem
	itemInfo := parentItem.GetItemBaseInfo()
	copier.Copy(&resParentItem, &itemInfo)
	sudItemIds := parentItem.GetSubItemIds()
	for _, id := range sudItemIds {
		resParentItem.SubItemIds = append(resParentItem.SubItemIds, int64(id))
	}
	res := &prpc.QuerySubItemsRes{
		ParentItem:    &resParentItem,
		TotalRowCount: int32(len(sudItemIds)),
	}

	for _, item := range items {
		var resItem prpc.CategoryItem
		itemInfo := item.GetItemBaseInfo()
		copier.Copy(&resItem, &itemInfo)
		sudItemIds := item.GetSubItemIds()
		for _, id := range sudItemIds {
			resItem.SubItemIds = append(resItem.SubItemIds, int64(id))
		}
		res.Items = append(res.Items, &resItem)
	}
	return res, nil
}

func (ser *CoreService) QueryItemInfo(ctx context.Context, req *prpc.QueryItemInfoReq) (*prpc.QueryItemInfoRes, error) {
	var userId ptype.UserID
	if len(req.ShareId) > 0 {
		si, err := ser.shares.GetShareItemInfo(req.ShareId)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "not found item")
		}
		if !ser.um.CategoryService().IsRelationOf(ptype.CategoryID(req.ItemId), si.ShareItemInfo.ItemId) {
			return nil, status.Error(codes.PermissionDenied, "not found item")
		}
		userId = si.UserId
	} else {
		ses := ser.getSession(ctx)
		if ses == nil {
			return nil, status.Error(codes.PermissionDenied, "not found session")
		}
		userId = ses.UserId
	}

	item, err := ser.um.CategoryService().GetItem(userId, ptype.CategoryID(req.ItemId))
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "not found")
	}
	res := &prpc.QueryItemInfoRes{}
	res.ItemInfo = &prpc.CategoryItem{}
	itemInfo := item.GetItemBaseInfo()
	copier.Copy(&res.ItemInfo, &itemInfo)
	subItemIds := item.GetSubItemIds()
	if len(subItemIds) > 0 {
		res.ItemInfo.SubItemIds = make([]int64, 0, len(subItemIds))
		for _, id := range subItemIds {
			res.ItemInfo.SubItemIds = append(res.ItemInfo.SubItemIds, int64(id))
		}
	}
	itemType := item.GetType()
	switch itemType {
	case prpc.CategoryItem_Video:
		res.VideoInfo = &prpc.Video{}
		vid := item.GetVideoId()
		lookPath := setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%d", vid)
		res.VideoInfo.Id = int64(vid)
		res.VideoInfo.SubtitlePaths = utils.GetNotZeroFilesByFileExtension(lookPath, []string{".vtt"})
	}
	return res, nil
}

func (ser *CoreService) AddBtVideos(ctx context.Context, req *prpc.AddBtVideosReq) (*prpc.AddBtVideosRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	if len(req.FileIndexes) == 0 || req.InfoHash == nil || req.CategoryItemId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	err := ser.um.AddBtVideos(&user.AddBtVideosParams{
		UserId:         ses.UserId,
		CategoryItemId: ptype.CategoryID(req.CategoryItemId),
		InfoHash:       bt.TranInfoHash(req.InfoHash),
		FileIndexes:    req.FileIndexes,
	})
	if err != nil {
		log.Warnf("[user] %d add videos err: %v", ses.UserId, err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &prpc.AddBtVideosRes{}, nil
}

func (ser *CoreService) ShareItem(ctx context.Context, req *prpc.ShareItemReq) (*prpc.ShareItemRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	_, err := ser.um.CategoryService().GetItem(ses.UserId, ptype.CategoryID(req.ItemId))
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "not found")
	}
	shareid, err := ser.shares.ShareCategoryItem(&ShareCategoryItemParams{
		UserId:    ses.UserId,
		ItemId:    ptype.CategoryID(req.ItemId),
		MaxCount:  0,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "not found")
	}
	return &prpc.ShareItemRes{
		ItemId:  req.ItemId,
		ShareId: shareid,
	}, nil
}

func (ser *CoreService) QuerySharedItems(ctx context.Context, req *prpc.QuerySharedItemsReq) (*prpc.QuerySharedItemsRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	sis := ser.shares.GetUserSharedItemInfos(ses.UserId)
	res := &prpc.QuerySharedItemsRes{}
	for _, si := range sis {
		res.SharedItems = append(res.SharedItems, &prpc.SharedItem{
			ItemId:  int64(si.ShareItemInfo.ItemId),
			ShareId: si.ShareId,
		})
	}
	return res, nil
}

func (ser *CoreService) GetShareItemInfo(shareid string) (*ShareInfo, error) {
	return ser.shares.GetShareItemInfo(shareid)
}

func (ser *CoreService) DelSharedItem(ctx context.Context, req *prpc.DelSharedItemReq) (*prpc.DelSharedItemRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	sit, err := ser.shares.GetShareItemInfo(req.ShareId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "not found shared item")
	}
	if sit.UserId != ses.UserId {
		return nil, status.Error(codes.PermissionDenied, "")
	}
	err = ser.shares.DelShare(req.ShareId)
	if err != nil {
		return nil, err
	}
	return &prpc.DelSharedItemRes{}, nil
}

func (ser *CoreService) UploadSubtitle(ctx context.Context, req *prpc.UploadSubtitleReq) (*prpc.UploadSubtitleRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	if req == nil || len(req.Subtitles) == 0 {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	err := ser.um.UploadSubtitle(ses.UserId, req)
	if err != nil {
		return nil, err
	}
	return &prpc.UploadSubtitleRes{}, nil
}

func (ser *CoreService) JoinChatRoom(req *prpc.JoinChatRoomReq, stream prpc.UserService_JoinChatRoomServer) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "")
	}
	ses := ser.getSession(stream.Context())
	if ses == nil {
		return status.Error(codes.PermissionDenied, "")
	}

	roomKey := utils.GetItemRoomKey(req.Room)
	joinParams := chat.JoinParams{
		RoomKey:          roomKey,
		SessionId:        ses.Id,
		MaxCacheNum:      30,
		MaxCacheDuration: time.Second * 1,
		NeedRecent:       true,
		SendFunc: func(cms []*chat.ChatMessage) {
			var scms []*prpc.ChatMessage
			for _, cm := range cms {
				user, _ := ser.um.LoadUser(cm.UserId)
				smsg := &prpc.ChatMessage{
					UserId:   int64(cm.UserId),
					UserName: user.GetUserName(),
					SentTime: cm.SentTime.UnixMilli(),
					Msg:      cm.Msg,
				}
				scms = append(scms, smsg)
			}
			stream.Send(&prpc.JoinChatRoomRes{
				Room:     req.Room,
				ChatMsgs: scms,
			})
		},
	}
	if req.Room.Type == prpc.Room_Danmaku {
		joinParams.NeedRecent = false
	}

	id, err := ser.rooms.Join(&joinParams)
	if err != nil {
		ser.rooms.CreateRoom(&chat.CreateRoomParams{
			RoomKey:       roomKey,
			ImmediatePush: false,
			Interval:      time.Second * 1,
		})
		id, _ = ser.rooms.Join(&joinParams)
	}

	room := ser.rooms.GetRoom(roomKey)
	if room == nil {
		return status.Error(codes.Internal, "")
	}
	select {
	case <-stream.Context().Done():
	case <-room.Context().Done():
	}
	room.Leave(id)
	return nil
}

func (ser *CoreService) SendMsg2ChatRoom(ctx context.Context, req *prpc.SendMsg2ChatRoomReq) (*prpc.SendMsg2ChatRoomRes, error) {
	if req == nil || req.GetChatMsg() == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}

	roomKey := utils.GetItemRoomKey(req.GetRoom())
	ser.rooms.Broadcast(roomKey, &chat.ChatMessage{
		UserId:   ses.UserId,
		SentTime: time.Now(),
		Msg:      req.GetChatMsg().Msg,
	})

	return &prpc.SendMsg2ChatRoomRes{}, nil
}

func (ser *CoreService) AddMagnetCategory(ctx context.Context, req *prpc.AddMagnetCategoryReq) (*prpc.AddMagnetCategoryRsp, error) {
	if req == nil || req.CategoryName == "" {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	_, err := ser.um.AddMagnetCategory(&user.AddMagnetCategoryParams{
		ParentId:  ptype.CategoryID(req.ParentId),
		Name:      req.CategoryName,
		Introduce: req.Introduce,
		Creator:   ses.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &prpc.AddMagnetCategoryRsp{}, nil
}

func (ser *CoreService) AddMagnetUri(ctx context.Context, req *prpc.AddMagnetUriReq) (*prpc.AddMagnetUriRsp, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}

	err := ser.um.AddMagnetUri(&user.AddMagnetUriParams{
		CategoryId: ptype.CategoryID(req.CategoryId),
		Introduce:  req.Introduce,
		Creator:    ses.UserId,
		Uri:        req.MagnetUri,
	})
	if err != nil {
		return nil, err
	}

	return &prpc.AddMagnetUriRsp{}, nil
}

func (ser *CoreService) QueryMagnet(ctx context.Context, req *prpc.QueryMagnetReq) (*prpc.QueryMagnetRsp, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}

	item, err := ser.um.CategoryService().GetItem(ptype.AdminId, ptype.CategoryID(req.ParentId))
	if err != nil {
		return nil, err
	}
	sudItemIds := item.GetSubItemIds()
	var items []*category.CategoryItem
	var totalRows int
	if len(req.SearchCond) > 0 {
		if !ser.um.CategoryService().IsRelationOf(ptype.CategoryID(req.ParentId), ser.um.GetMagnetRootId()) {
			return nil, status.Error(codes.PermissionDenied, "not found parent id")
		}
		type searCond struct {
			ExistedWords    []string `json:"ExistedWords"`
			NotExistedWords []string `json:"NotExistedWords"`
		}
		var cond searCond
		err := json.Unmarshal([]byte(req.SearchCond), &cond)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "")
		}
		params := &category.SearchParams{
			Querier:         ptype.AdminId,
			RootId:          ptype.CategoryID(req.ParentId),
			ExistedWords:    cond.ExistedWords,
			NotExistedWords: cond.NotExistedWords,
			PageNum:         req.PageNum,
			Rows:            req.Rows,
		}
		totalRows, err = ser.um.CategoryService().SearchRows(params)
		if err != nil {
			return nil, err
		}
		items, err = ser.um.CategoryService().Search(params)
		if err != nil {
			return nil, err
		}
	} else {
		items, err = ser.um.QueryMagnetCategorys(&user.QueryCategoryParams{
			ParentId: ptype.CategoryID(req.ParentId),
			PageNum:  req.PageNum,
			Rows:     req.Rows,
		})
		totalRows = len(sudItemIds)
		if err != nil {
			return nil, err
		}
	}

	res := &prpc.QueryMagnetRsp{
		TotalRowCount: int32(totalRows),
	}
	var resItem prpc.CategoryItem
	itemInfo := item.GetItemBaseInfo()
	copier.Copy(&resItem, &itemInfo)

	for _, id := range sudItemIds {
		resItem.SubItemIds = append(resItem.SubItemIds, int64(id))
	}
	res.Items = append(res.Items, &resItem)
	for _, item := range items {
		var resItem prpc.CategoryItem
		itemInfo := item.GetItemBaseInfo()
		copier.Copy(&resItem, &itemInfo)
		sudItemIds := item.GetSubItemIds()
		for _, id := range sudItemIds {
			resItem.SubItemIds = append(resItem.SubItemIds, int64(id))
		}
		res.Items = append(res.Items, &resItem)
	}
	return res, nil
}

func (ser *CoreService) DelMagnetCategory(ctx context.Context, req *prpc.DelMagnetCategoryReq) (*prpc.DelMagnetCategoryRsp, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}

	err := ser.um.DelMagnetCategory(ses.UserId, ptype.CategoryID(req.Id))
	if err != nil {
		return nil, err
	}
	return &prpc.DelMagnetCategoryRsp{}, nil
}

func (ser *CoreService) GetBtMeta(ctx context.Context, req *prpc.GetBtMetaReq) (*prpc.GetBtMetaRsp, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	if req == nil || req.Req == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	req.Req.StopAfterGotMeta = true
	_, err := ser.um.Download(&bt.DownloadParams{
		UserId: ses.UserId,
		Req:    req.Req,
	})
	if err != nil {
		return nil, err
	}
	return &prpc.GetBtMetaRsp{}, nil
}

func (ser *CoreService) NewBtHlsTask(ctx context.Context, req *prpc.NewBtHlsTaskReq) (*prpc.NewBtHlsTaskRsp, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	if req == nil || req.Req == nil || req.CategoryParentId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	err := ser.um.GetTasks().NewBtHlsTask(&task.NewBtHlsParams{
		UserId:           ses.UserId,
		ParentId:         ptype.CategoryID(req.CategoryParentId),
		Bt:               &ser.um,
		DownloadReq:      req.Req,
		CategorySer:      ser.um.CategoryService(),
		RecursiveNewPath: req.RecursiveNewPath,
	})
	if err != nil {
		return nil, err
	}
	return &prpc.NewBtHlsTaskRsp{}, nil
}

func (ser *CoreService) RenameBtVideoName(ctx context.Context, req *prpc.RenameBtVideoNameReq) (*prpc.RenameBtVideoNameRsp, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	err := ser.um.RenameBtVideoName(&user.RenameBtVideoNameParams{
		Who:      ses.UserId,
		ParentId: ptype.CategoryID(req.ItemId),
		RefName:  req.RefName,
	})
	if err != nil {
		return nil, err
	}
	return &prpc.RenameBtVideoNameRsp{}, nil
}
