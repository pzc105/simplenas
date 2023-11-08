package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"pnas/bt"
	"pnas/category"
	"pnas/crawler"
	"pnas/log"
	"pnas/phttp"
	"pnas/prpc"
	"pnas/service/chat"
	"pnas/service/session"
	"pnas/setting"
	"pnas/user"
	"pnas/utils"
	"pnas/video"
	"sort"
	"strconv"
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
	notCheckTokenMethods []string

	sessions session.ISessions

	bt              bt.BtClient
	btStatusPushMtx sync.Mutex
	btStatusPush    map[int64]chan *prpc.StatusRespone

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
	ser.btStatusPush = make(map[int64]chan *prpc.StatusRespone)

	ser.notCheckTokenMethods = []string{"Register", "IsUsedEmail", "Login", "FastLogin", "QueryItemInfo", "QuerySubItems"}
	sort.Strings(ser.notCheckTokenMethods)
	ser.bt.Init(bt.WithOnStatus(ser.handleBtStatus),
		bt.WithOnTorrentInfo(ser.handleTorrentInfo),
		bt.WithOnConnect(ser.handleBtClientConnected),
		bt.WithOnFileCompleted(ser.handleBtFileCompleted))
	ser.um.Init()

	go crawler.Go36dmBackgroup(&ser.um, &ser.bt)

	var rooms chat.Rooms
	rooms.Init()
	ser.rooms = &rooms

	sm := &ShareManager{}
	sm.Init()
	ser.shares = sm
}

func (ser *CoreService) Serve() {
	creds, err := credentials.NewServerTLSFromFile(setting.GS().Server.CrtFile, setting.GS().Server.KeyFile)
	if err != nil {
		log.Panic(err)
	}

	ser.grpcSer = grpc.NewServer(grpc.Creds(creds),
		grpc.UnaryInterceptor(ser.CheckToken),
		grpc.StreamInterceptor(ser.StreamCheckToken))
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

func (ser *CoreService) recvDanmaku(vid video.ID, danmaku string) {
	roomKey := getItemRoomKey(&prpc.Room{
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
	ser.bt.Close()
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

func (ser *CoreService) handleBtClientConnected() {
	log.Info("connected to bt service")
	resumeData := user.LoadDownloadingTorrent()
	for _, resume := range resumeData {
		var req prpc.DownloadRequest
		req.Type = prpc.DownloadRequest_Resume
		req.Content = resume
		req.SavePath = setting.GS().Bt.SavePath
		_, err := ser.bt.Download(context.Background(), &req)
		if err != nil {
			log.Warnf("[bt] failed to download, err: %v", err)
		}
	}
}

func (ser *CoreService) handleTorrentInfo(tis *prpc.TorrentInfoRes) {
	var b bt.TorrentBase
	ti := tis.GetTi()
	copier.Copy(&b, ti)
	var files []bt.File
	copier.Copy(&files, ti.Files)
	go ser.um.UpdateTorrent(&user.UpdateTorrentParams{
		Base:       &b,
		State:      ti.GetState(),
		FileNames:  files,
		ResumeData: ti.GetResumeData(),
	})
}

func (ser *CoreService) handleBtStatus(sr *prpc.StatusRespone) {
	ser.btStatusPushMtx.Lock()
	defer ser.btStatusPushMtx.Unlock()
	for sid, ch := range ser.btStatusPush {
		ses, err := ser.sessions.GetSession3(sid)
		if err != nil {
			delete(ser.btStatusPush, sid)
			continue
		}
		var r prpc.StatusRespone
		for _, st := range sr.StatusArray {
			tInfoHash := TranInfoHash(st.InfoHash)
			if !ser.um.HasTorrent(ses.UserId, tInfoHash) {
				continue
			}
			r.StatusArray = append(r.StatusArray, st)
		}
		if len(r.StatusArray) > 0 && len(ch) == 0 {
			ch <- &r
		}
	}
}

func (ser *CoreService) handleBtFileCompleted(fs *prpc.FileCompletedRes) {
	lfc := &user.FileCompleted{
		InfoHash:  TranInfoHash(fs.InfoHash),
		FileIndex: fs.FileIndex,
	}
	go ser.um.BtFileStateComplete(lfc)
}

func (ser *CoreService) getSession(ctx context.Context) *session.Session {
	s, _ := ser.sessions.GetSession2(ctx)
	return s
}

func (ser *CoreService) CheckToken(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	methodName := info.FullMethod[len(prpc.UserService_ServiceDesc.ServiceName)+2:]
	i := sort.SearchStrings(ser.notCheckTokenMethods, methodName)
	if i == len(ser.notCheckTokenMethods) || ser.notCheckTokenMethods[i] != methodName {
		_, err := ser.sessions.GetSession2(ctx)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "")
		}
	}
	return handler(ctx, req)
}

func (ser *CoreService) StreamCheckToken(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	if srv == nil {
		return status.Error(codes.InvalidArgument, "")
	}
	methodName := info.FullMethod[len(prpc.UserService_ServiceDesc.ServiceName)+2:]
	for _, m := range ser.notCheckTokenMethods {
		if m == methodName {
			continue
		}
		_, err := ser.sessions.GetSession2(ss.Context())
		if err != nil {
			return status.Error(codes.InvalidArgument, "")
		}
		break
	}
	return handler(srv, ss)
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
		Email:  userInfo.GetEmail(),
		Name:   userInfo.GetName(),
		Passwd: userInfo.GetPasswd(),
		Auth:   utils.NewBitSet(user.AuthMax, user.AuthAdmin),
	}
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

	if user.IsUsedEmail(emailInfo.GetEmail()) {
		return nil, status.Error(codes.AlreadyExists, "existed email")
	}
	return &prpc.IsUsedEmailRet{}, nil
}

func (ser *CoreService) Login(
	ctx context.Context,
	loginInfo *prpc.LoginInfo) (*prpc.LoginRet, error) {

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

	res, err := ser.bt.Parse(context.Background(), req)
	if err != nil {
		return nil, err
	}
	resumeData, err := user.LoadTorrent(TranInfoHash(res.InfoHash))
	if err == nil {
		req.Type = prpc.DownloadRequest_Resume
		req.Content = resumeData
	}

	req.SavePath = setting.GS().Bt.SavePath
	res, err = ser.bt.Download(context.Background(), req)
	if err == nil {
		log.Infof("[bt] user %d token:%s add torrent: %s",
			ses.UserId,
			ses.Token,
			hex.EncodeToString(res.InfoHash.Hash))
		ser.um.AddTorrent(ses.UserId, &bt.TorrentBase{
			InfoHash: TranInfoHash(res.InfoHash),
		})
		return res, nil
	} else {
		return res, status.Error(codes.InvalidArgument, "")
	}
}

func (ser *CoreService) RemoveTorrent(ctx context.Context, req *prpc.RemoveTorrentReq) (*prpc.RemoveTorrentRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "")
	}

	if ses.UserId == user.AdminId {
		err := ser.um.RemoveTorrent(user.AdminId, TranInfoHash(req.InfoHash))
		if err != nil {
			log.Warnf("[user] user %d failed to remove torrent %s, err: %v", ses.UserId, hex.EncodeToString(req.InfoHash.Hash), err)
			return nil, status.Error(codes.Unknown, "")
		}
		_, err = ser.bt.RemoveTorrent(context.Background(), req)
		if err != nil {
			log.Warnf("[bt] user %d failed to remove torrent %s, err: %v", ses.UserId, hex.EncodeToString(req.InfoHash.Hash), err)
			return nil, status.Error(codes.Unknown, "")
		}
		return &prpc.RemoveTorrentRes{}, nil
	}

	err := ser.um.RemoveUserTorrent(ses.UserId, TranInfoHash(req.InfoHash))
	if err != nil {
		log.Warnf("[user] user %d failed to remove torrent %s, err: %v", ses.UserId, hex.EncodeToString(req.InfoHash.Hash), err)
		return nil, status.Error(codes.Unknown, "")
	}
	return &prpc.RemoveTorrentRes{}, nil
}

func (ser *CoreService) OnStatus(statusReq *prpc.StatusRequest, stream prpc.UserService_OnStatusServer) error {
	ses := ser.getSession(stream.Context())
	if ses == nil {
		return status.Error(codes.PermissionDenied, "")
	}

	ser.btStatusPushMtx.Lock()
	ch, ok := ser.btStatusPush[ses.Id]
	if ok {
		close(ch)
	}
	ch = make(chan *prpc.StatusRespone, 1)
	ser.btStatusPush[ses.Id] = ch
	ser.btStatusPushMtx.Unlock()

	for {
		b := false
		select {
		case sr, ok := <-ch:
			if !ok {
				return nil
			}
			stream.Send(sr)
		case <-stream.Context().Done():
			b = true
		}
		if b {
			break
		}
	}
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
	ms, err := ser.um.QueryBtVideoMetadata(ses.UserId, TranInfoHash(req.InfoHash))
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
		ParentId:     category.ID(req.ParentId),
		Creator:      int64(ses.UserId),
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
	err := ser.um.DelCategoryItem(ses.UserId, category.ID(req.GetItemId()))
	if err != nil {
		log.Warnf("[user] %d delete category %d err: %v", ses.UserId, req.GetItemId(), err)
		return nil, status.Error(codes.Unknown, "")
	}
	return &prpc.DelCategoryItemRes{}, nil
}

func (ser *CoreService) QuerySubItems(ctx context.Context, req *prpc.QuerySubItemsReq) (*prpc.QuerySubItemsRes, error) {
	var userId user.ID
	if len(req.ShareId) > 0 {
		si, err := ser.shares.GetShareItemInfo(req.ShareId)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "not found item")
		}
		if !ser.um.IsRelationOf(category.ID(req.ParentId), si.ShareItemInfo.ItemId) {
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

	parentItem, err := ser.um.CategoryService().GetItem(int64(userId), category.ID(req.ParentId))
	if err != nil {
		log.Warnf("[user] %d query category %d err: %v", userId, req.ParentId, err)
		return nil, status.Error(codes.PermissionDenied, "not found")
	}
	items, err := ser.um.CategoryService().GetItemsByParent(
		&category.GetItemsByParentParams{
			Querier:  int64(userId),
			ParentId: category.ID(req.ParentId),
			PageNum:  req.PageNum,
			Rows:     req.Rows,
		})
	if err != nil {
		return nil, err
	}

	var resParentItem prpc.CategoryItem
	itemInfo := parentItem.GetItemInfo()
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
		itemInfo := item.GetItemInfo()
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
	var userId user.ID
	if len(req.ShareId) > 0 {
		si, err := ser.shares.GetShareItemInfo(req.ShareId)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "not found item")
		}
		if !ser.um.IsRelationOf(category.ID(req.ItemId), si.ShareItemInfo.ItemId) {
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

	item, err := ser.um.CategoryService().GetItem(int64(userId), category.ID(req.ItemId))
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "not found")
	}
	res := &prpc.QueryItemInfoRes{}
	res.ItemInfo = &prpc.CategoryItem{}
	itemInfo := item.GetItemInfo()
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
	if len(req.FileIndexes) == 0 || req.InfoHash == nil {
		return nil, status.Error(codes.InvalidArgument, "")
	}
	err := ser.um.AddBtVideos(&user.AddBtVideosParams{
		UserId:         ses.UserId,
		CategoryItemId: category.ID(req.CategoryItemId),
		InfoHash:       TranInfoHash(req.InfoHash),
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
	_, err := ser.um.CategoryService().GetItem(int64(ses.UserId), category.ID(req.ItemId))
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "not found")
	}
	shareid, err := ser.shares.ShareCategoryItem(&ShareCategoryItemParams{
		UserId:    ses.UserId,
		ItemId:    category.ID(req.ItemId),
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

func (ser *CoreService) RefreshSubtitle(ctx context.Context, req *prpc.RefreshSubtitleReq) (*prpc.RefreshSubtitleRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	item, err := ser.um.CategoryService().GetItem(int64(ses.UserId), category.ID(req.ItemId))
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "not found")
	}
	vid, _ := strconv.ParseInt(item.GetItemInfo().ResourcePath, 10, 64)
	err = ser.um.RefreshSubtitle(video.ID(vid))
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &prpc.RefreshSubtitleRes{}, nil
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

	roomKey := getItemRoomKey(req.Room)
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

	roomKey := getItemRoomKey(req.GetRoom())
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
		ParentId:  category.ID(req.ParentId),
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
		CategoryId: category.ID(req.CategoryId),
		Uri:        req.MagnetUri,
		Introduce:  req.Introduce,
		Creator:    ses.UserId,
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

	item, err := ser.um.CategoryService().GetItem(category.AdminId, category.ID(req.ParentId))
	if err != nil {
		return nil, err
	}
	sudItemIds := item.GetSubItemIds()
	var items []*category.CategoryItem
	var totalRows int
	if len(req.SearchCond) > 0 {
		if !ser.um.CategoryService().IsRelationOf(category.ID(req.ParentId), ser.um.GetMagnetRootId()) {
			return nil, status.Error(codes.PermissionDenied, "not found parent id")
		}
		params := &category.SearchParams{
			Querier:      category.AdminId,
			RootId:       category.ID(req.ParentId),
			ExistedWords: req.SearchCond,
			PageNum:      req.PageNum,
			Rows:         req.Rows,
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
			ParentId: category.ID(req.ParentId),
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
	itemInfo := item.GetItemInfo()
	copier.Copy(&resItem, &itemInfo)

	for _, id := range sudItemIds {
		resItem.SubItemIds = append(resItem.SubItemIds, int64(id))
	}
	res.Items = append(res.Items, &resItem)
	for _, item := range items {
		var resItem prpc.CategoryItem
		itemInfo := item.GetItemInfo()
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

	err := ser.um.DelMagnetCategory(ses.UserId, category.ID(req.Id))
	if err != nil {
		return nil, err
	}
	return &prpc.DelMagnetCategoryRsp{}, nil
}
