package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"pnas/bt"
	"pnas/category"
	"pnas/log"
	"pnas/phttp"
	"pnas/prpc"
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

type CoreServiceInterface interface {
	GetSession(*http.Request) *session
	GetUserManager() *user.UserManger
	GetShareItemInfo(shareid string) (*ShareItemInfo, error)
}

type CoreService struct {
	CoreServiceInterface
	prpc.UnimplementedUserServiceServer

	sessionsMtx sync.Mutex
	sessions    map[int64]*session

	notCheckTokenMethods []string

	bt bt.BtClient

	um          user.UserManger
	shutDownCtx context.Context
	closeFunc   context.CancelFunc

	wg sync.WaitGroup

	grpcSer   *grpc.Server
	httpSer   *http.Server
	videoSer  *VideoService
	posterSer *PosterService

	shares ShareManager
}

func (ser *CoreService) Init() {
	InitIdPool()
	ser.notCheckTokenMethods = []string{"Register", "IsUsedEmail", "Login", "FastLogin", "QueryItemInfo", "QuerySubItems"}
	sort.Strings(ser.notCheckTokenMethods)
	ser.sessions = make(map[int64]*session)
	ser.bt.Init(bt.WithOnStatus(ser.handleBtStatus),
		bt.WithOnTorrentInfo(ser.handleTorrentInfo),
		bt.WithOnConnect(ser.handleBtClientConnected),
		bt.WithOnFileCompleted(ser.handleBtFileCompleted))
	ser.um.Init()
	ser.shares.Init()
}

func (ser *CoreService) Serve() {
	creds, err := credentials.NewServerTLSFromFile(setting.GS.Server.CrtFile, setting.GS.Server.KeyFile)
	if err != nil {
		log.Panic(err)
	}

	ser.grpcSer = grpc.NewServer(grpc.Creds(creds),
		grpc.UnaryInterceptor(ser.CheckToken),
		grpc.StreamInterceptor(ser.StreamCheckToken))
	prpc.RegisterUserServiceServer(ser.grpcSer, ser)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d",
		setting.GS.Server.BoundIp,
		setting.GS.Server.Port))
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
	ser.videoSer = newVideoService(ser, router.PathPrefix("/video").Subrouter())
	ser.posterSer = newPosterService(ser, router.PathPrefix("/poster").Subrouter())

	ser.httpSer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", setting.GS.Server.BoundIp, setting.GS.Server.WebPort),
		Handler: router,
	}

	ser.shutDownCtx, ser.closeFunc = context.WithCancel(context.Background())

	if err := ser.httpSer.ListenAndServeTLS(setting.GS.Server.CrtFile, setting.GS.Server.KeyFile); err != nil {
		log.Panic(err)
	}
}

func (ser *CoreService) Close() {
	ser.closeFunc()
	ser.bt.Close()
	ser.httpSer.Shutdown(context.Background())
	ser.grpcSer.GracefulStop()
	ser.wg.Wait()
}

func (ser *CoreService) Wait() {
	// TODO this is not right
	ser.wg.Wait()
}

func (ser *CoreService) GetUserManager() *user.UserManger {
	return &ser.um
}

func (ser *CoreService) GetSession(r *http.Request) *session {
	_, cId := GetTokenAndIdByCookie(r.Header.Get("cookie"))
	ser.sessionsMtx.Lock()
	defer ser.sessionsMtx.Unlock()
	s, ok := ser.sessions[cId]
	if !ok {
		return nil
	}
	return s
}

func (ser *CoreService) handleBtClientConnected() {
	log.Info("connected to bt service")
	resumeData := ser.um.LoadDownloadingTorrent()
	for _, resume := range resumeData {
		var req prpc.DownloadRequest
		req.Type = prpc.DownloadRequest_Resume
		req.Content = resume
		req.SavePath = setting.GS.Bt.SavePath
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
	go ser.um.UpdateTorrent(&b, ti.GetState(), files, ti.GetResumeData())
}

func (ser *CoreService) handleBtStatus(sr *prpc.StatusRespone) {
	ser.sessionsMtx.Lock()
	defer ser.sessionsMtx.Unlock()

	for _, ses := range ser.sessions {
		if !ses.needPush.Load() {
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
		if len(r.StatusArray) > 0 {
			ses.btStatusCh <- &r
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

func (user *CoreService) verifySession(ctx context.Context) (ok bool) {
	cToken, cId := GetTokenAndId(ctx)
	user.sessionsMtx.Lock()
	defer user.sessionsMtx.Unlock()
	eSession, ok := user.sessions[cId]
	if !ok || cToken != eSession.Token {
		return false
	}
	return true
}

func (ser *CoreService) verifyToen(ctx context.Context) (ok bool, sess *session, cId int64) {
	cToken, cId := GetTokenAndId(ctx)
	ser.sessionsMtx.Lock()
	defer ser.sessionsMtx.Unlock()
	eSession, ok := ser.sessions[cId]
	if !ok {
		rSession, err := loadSession(cId)
		if err != nil || cToken != rSession.Token {
			return false, nil, -1
		}
		return true, rSession, cId
	} else if cToken != eSession.Token {
		return false, nil, -1
	}
	return true, eSession, cId
}

func (ser *CoreService) getSession(ctx context.Context) *session {
	_, cId := GetTokenAndId(ctx)
	ser.sessionsMtx.Lock()
	defer ser.sessionsMtx.Unlock()
	eSession, ok := ser.sessions[cId]
	if ok {
		return eSession
	}
	return nil
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
		ok := ser.verifySession(ctx)
		if !ok {
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
		ok := ser.verifySession(ss.Context())
		if !ok {
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

	session := ser.getSession(ctx)
	if session == nil {
		var md metadata.MD
		md, session = GenCookieSession(loginInfo.RememberMe, userInfo.Id)
		grpc.SendHeader(ctx, md)
	} else {
		var md metadata.MD
		md, session = GenCookieSessionById(session.Id, loginInfo.RememberMe, userInfo.Id)
		grpc.SendHeader(ctx, md)
	}

	ser.sessionsMtx.Lock()
	defer ser.sessionsMtx.Unlock()
	ser.sessions[session.Id] = session

	if loginInfo.RememberMe {
		e := saveSession(session)
		if e != nil {
			log.Warnf("[user] %d failed to save session, err: %v", userInfo.Id, e)
		}
	}

	return &prpc.LoginRet{
		Token: session.Token,
		UserInfo: &prpc.UserInfo{
			Id:              int64(userInfo.Id),
			Name:            userInfo.Name,
			Email:           userInfo.Email,
			HomeDirectoryId: int64(userInfo.HomeDirectoryId),
		},
		RememberMe: loginInfo.RememberMe,
	}, nil
}

func (ser *CoreService) FastLogin(
	ctx context.Context,
	loginInfo *prpc.LoginInfo) (*prpc.LoginRet, error) {

	ok, oldSession, cId := ser.verifyToen(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "")
	}

	user, err := ser.um.LoadUser(oldSession.UserId)
	if err != nil {
		log.Warnf("[user] %d failed to fast login, err: %v", user.GetUserInfo().Id, err)
		return nil, status.Error(codes.NotFound, "")
	}

	md, session := GenCookieSessionById(cId, loginInfo.RememberMe, oldSession.UserId)
	grpc.SendHeader(ctx, md)
	ser.sessionsMtx.Lock()
	defer ser.sessionsMtx.Unlock()
	ser.sessions[cId] = session

	if loginInfo.RememberMe {
		saveSession(session)
	}

	userInfo := user.GetUserInfo()
	return &prpc.LoginRet{
		Token: session.Token,
		UserInfo: &prpc.UserInfo{
			Id:              int64(userInfo.Id),
			Name:            userInfo.Name,
			Email:           userInfo.Email,
			HomeDirectoryId: int64(userInfo.HomeDirectoryId),
		},
		RememberMe: loginInfo.RememberMe,
	}, nil
}

func (ser *CoreService) IsLogined(context.Context, *prpc.LoginInfo) (*prpc.LoginRet, error) {
	return &prpc.LoginRet{}, nil
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
	resumeData, err := ser.um.LoadTorrent(TranInfoHash(res.InfoHash))
	if err == nil {
		req.Type = prpc.DownloadRequest_Resume
		req.Content = resumeData
	}

	req.SavePath = setting.GS.Bt.SavePath
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

	if ses.needPush.Load() {
		return status.Error(codes.InvalidArgument, "")
	}
	ses.btStatusCh = make(chan *prpc.StatusRespone)
	ses.needPush.Store(true)
	defer ses.needPush.Store(false)
	for {
		b := false
		select {
		case sr, ok := <-ses.btStatusCh:
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
		sii, err := ser.shares.GetShareItemInfo(req.ShareId)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "not found item")
		}
		if !ser.um.IsItemShared(sii.ItemId, category.ID(req.ParentId)) {
			return nil, status.Error(codes.PermissionDenied, "not found item")
		}
		userId = sii.UserId
	} else {
		ses := ser.getSession(ctx)
		if ses == nil {
			return nil, status.Error(codes.PermissionDenied, "not found session")
		}
		userId = ses.UserId
	}

	parentItem, err := ser.um.QueryItem(userId, category.ID(req.ParentId))
	if err != nil {
		log.Warnf("[user] %d query category %d err: %v", userId, req.ParentId, err)
		return nil, status.Error(codes.PermissionDenied, "not found")
	}
	items := ser.um.QueryItems(userId, category.ID(req.ParentId))

	var resParentItem prpc.CategoryItem
	itemInfo := parentItem.GetItemInfo()
	copier.Copy(&resParentItem, &itemInfo)
	sudItemIds := parentItem.GetSubItemIds()
	for _, id := range sudItemIds {
		resParentItem.SubItemIds = append(resParentItem.SubItemIds, int64(id))
	}
	res := &prpc.QuerySubItemsRes{
		ParentItem: &resParentItem,
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
		sii, err := ser.shares.GetShareItemInfo(req.ShareId)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "not found item")
		}
		if !ser.um.IsItemShared(sii.ItemId, category.ID(req.ItemId)) {
			return nil, status.Error(codes.PermissionDenied, "not found item")
		}
		userId = sii.UserId
	} else {
		ses := ser.getSession(ctx)
		if ses == nil {
			return nil, status.Error(codes.PermissionDenied, "not found session")
		}
		userId = ses.UserId
	}

	item, err := ser.um.QueryItem(userId, category.ID(req.ItemId))
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
		vid, _ := strconv.ParseInt(itemInfo.ResourcePath, 10, 64)
		lookPath := setting.GS.Server.HlsPath + fmt.Sprintf("/vid_%d", vid)
		res.VideoInfo.Id = vid
		res.VideoInfo.SubtitlePaths = utils.GetFilesByFileExtension(lookPath, []string{".vtt"})
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
	_, err := ser.um.QueryItem(ses.UserId, category.ID(req.ItemId))
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

func (ser *CoreService) GetShareItemInfo(shareid string) (*ShareItemInfo, error) {
	return ser.shares.GetShareItemInfo(shareid)
}

func (ser *CoreService) RefreshSubtitle(ctx context.Context, req *prpc.RefreshSubtitleReq) (*prpc.RefreshSubtitleRes, error) {
	ses := ser.getSession(ctx)
	if ses == nil {
		return nil, status.Error(codes.PermissionDenied, "not found session")
	}
	item, err := ser.um.QueryItem(ses.UserId, category.ID(req.ItemId))
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
