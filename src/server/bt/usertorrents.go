package bt

import (
	"context"
	"encoding/hex"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/setting"
	"pnas/utils"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AdminId = 1
)

type FileCompleted struct {
	InfoHash  InfoHash
	FileIndex int32
}

type UserTorrentsImpl struct {
	mtx      sync.Mutex
	torrents map[InfoHash]*Torrent
	users    map[ptype.UserID]*userData

	shutDownCtx context.Context
	closeFunc   context.CancelFunc
	wg          sync.WaitGroup

	btClient BtClient

	trackerSub trackerSub

	callbackTaskQueue utils.TaskQueue
}

func (ut *UserTorrentsImpl) Init() {
	ut.shutDownCtx, ut.closeFunc = context.WithCancel(context.Background())

	ut.torrents = make(map[InfoHash]*Torrent)
	ut.users = make(map[ptype.UserID]*userData)
	ut.load()

	ut.btClient.Init(WithOnStatus(ut.onBtStatus),
		WithOnConnect(ut.handleBtClientConnected),
		WithOnFileCompleted(ut.handleBtFileCompleted))

	ut.callbackTaskQueue.Init()

	ut.trackerSub.Init()

	ut.wg.Add(1)
	go func() {
		defer ut.wg.Add(-1)
		ticker := time.NewTicker(time.Minute * 1)
	loop:
		for {
			select {
			case <-ticker.C:
				rsp, err := ut.btClient.GetSessionParams(ut.shutDownCtx, &prpc.GetSessionParamsReq{})
				if err == nil {
					saveBtSessionParams(rsp.ResumeData)
				}
			case <-ut.shutDownCtx.Done():
				break loop
			}
		}
	}()
}

func (ut *UserTorrentsImpl) load() {
retry:
	sql := `select user_id, torrent_id from user_torrent`
	rows, err := db.Query(sql)
	if err != nil {
		log.Warnf("load user torrent err: %v", err)
		<-time.After(time.Second * 5)
		goto retry
	}
	defer rows.Close()
	for rows.Next() {
		var uid ptype.UserID
		var tid ptype.TorrentID
		err := rows.Scan(&uid, &tid)
		if err != nil {
			log.Warnf("load user torrent err: %v", err)
			return
		}
		var t *Torrent
		var ok bool
		ut.mtx.Lock()
		infoHash, err := loadInfoHash(tid)
		if err != nil {
			ut.mtx.Unlock()
			continue
		}
		t, ok = ut.torrents[*infoHash]
		if !ok {
			t = loadTorrent(&ut.btClient, tid)
			if t == nil {
				ut.mtx.Unlock()
				continue
			}
			ut.torrents[t.base.InfoHash] = t
		}
		ud := ut.getUserDataLocked(uid)
		ut.mtx.Unlock()
		ud.initTorrent(t)
		t.addUser(ud)
	}
}

func (ut *UserTorrentsImpl) Close() {
	ut.btClient.Close()
	ut.callbackTaskQueue.Close(utils.CloseWayDrained)
	ut.wg.Wait()
}

func (ut *UserTorrentsImpl) handleBtClientConnected() {
	log.Info("connected to bt service")

	rsp, err := ut.btClient.InitedSession(context.Background(), &prpc.InitedSessionReq{})
	if err != nil {
		return
	}

	if !rsp.Inited {
		resume, _ := loadBtSessionParams()
		req := &prpc.InitSessionReq{
			ProxyHost:         setting.GS().Bt.ProxyHost,
			ProxyPort:         setting.GS().Bt.ProxyPort,
			ProxyType:         setting.GS().Bt.ProxyType,
			UploadRateLimit:   setting.GS().Bt.UploadRateLimit,
			DownloadRateLimit: setting.GS().Bt.DownloadRateLimit,
			HashingThreads:    setting.GS().Bt.HashingThreads,
			ListenInterfaces:  setting.GS().Bt.ListenInterfaces,
			ResumeData:        resume,
		}
		ut.btClient.InitSession(context.Background(), req)
	}

	ut.mtx.Lock()
	defer ut.mtx.Unlock()

	for _, t := range ut.torrents {
		req := &prpc.DownloadRequest{}
		resume, err := loadResumeData(&t.base.InfoHash)
		if err == nil {
			req.Type = prpc.DownloadRequest_Resume
			req.Content = resume
		} else {
			uri, err := getMagnetByInfoHash(&t.base.InfoHash)
			if err != nil {
				continue
			}
			req.Type = prpc.DownloadRequest_MagnetUri
			req.Content = []byte(uri)
		}
		_, err = ut.download(&t.base.InfoHash, req)
		if err != nil {
			log.Warnf("failed download err: %v", err)
		}
	}
}

func (ut *UserTorrentsImpl) download(infoHash *InfoHash, req *prpc.DownloadRequest) (*prpc.DownloadRespone, error) {
	req.SavePath = setting.GS().Bt.SavePath
	if req.Type != prpc.DownloadRequest_Resume {
		resumeData, err := loadResumeData(infoHash)
		if err == nil {
			req.Type = prpc.DownloadRequest_Resume
			req.Content = resumeData
		}
	}
	trackersstr := setting.GS().Bt.Trackers
	if len(trackersstr) > 0 {
		trackers := strings.Split(trackersstr, ",")
		for i := range trackers {
			req.Trackers = append(req.Trackers, strings.Trim(trackers[i], " "))
		}
	}
	trackers := ut.trackerSub.GetTrackers()
	req.Trackers = append(req.Trackers, trackers...)
	return ut.btClient.Download(context.Background(), req)
}

func (ut *UserTorrentsImpl) SetTaskCallback(params *SetTaskCallbackParams) {
	ut.getUserData(params.UserId).setTaskCallback(params)
}

func (ut *UserTorrentsImpl) SetSessionCallback(userId ptype.UserID, sid ptype.SessionID, callback UserOnBtStatusCallback) {
	ut.getUserData(userId).setSessionCallback(sid, callback)
}

func (ut *UserTorrentsImpl) onBtStatus(sr *prpc.BtStatusRespone) {
	for _, s := range sr.StatusArray {
		ut.handleBtStatus(s)
	}
}

func (ut *UserTorrentsImpl) handleBtStatus(s *prpc.TorrentStatus) {
	ut.mtx.Lock()
	t, ok := ut.torrents[*TranInfoHash(s.InfoHash)]
	if !ok {
		ut.mtx.Unlock()
		return
	}
	ut.mtx.Unlock()
	if !t.hasBaseInfo() && s.State != prpc.BtStateEnum_downloading_metadata {
		req := &prpc.GetTorrentInfoReq{
			InfoHash: s.InfoHash,
		}
		tRes, err := ut.btClient.GetTorrentInfo(context.Background(), req)
		if len(t.GetMagnetUri()) == 0 {
			mRes, err := ut.btClient.GetMagnetUri(context.Background(), &prpc.GetMagnetUriReq{
				Type:     prpc.GetMagnetUriReq_InfoHash,
				InfoHash: s.InfoHash,
			})
			if err == nil && len(mRes.MagnetUri) > 0 {
				t.UpdateMagnetUri(mRes.MagnetUri)
			}
		}
		if err != nil {
			goto updateSt
		}
		t.updateTorrentInfo(tRes.TorrentInfo)
	}
updateSt:
	t.updateStatus(s)
}

func (ut *UserTorrentsImpl) handleBtFileCompleted(fs *prpc.FileCompletedRes) {
	lfc := &FileCompleted{
		InfoHash:  *TranInfoHash(fs.InfoHash),
		FileIndex: fs.FileIndex,
	}
	go ut.btFileStateComplete(lfc)
}

func (ut *UserTorrentsImpl) btFileStateComplete(fs *FileCompleted) {
	ut.mtx.Lock()
	t, ok := ut.torrents[fs.InfoHash]
	ut.mtx.Unlock()
	if !ok {
		return
	}
	t.UpdateFileState(int(fs.FileIndex), prpc.BtFile_completed)
}

func (ut *UserTorrentsImpl) getUserDataLocked(uid ptype.UserID) *userData {
	ud, ok := ut.users[uid]
	if ok {
		return ud
	}
	ud = &userData{
		userId: uid,
	}
	ud.init(&ut.callbackTaskQueue)
	ut.users[uid] = ud
	return ud
}

func (ut *UserTorrentsImpl) getUserData(uid ptype.UserID) *userData {
	ut.mtx.Lock()
	defer ut.mtx.Unlock()
	return ut.getUserDataLocked(uid)
}

func (ut *UserTorrentsImpl) NewTorrentByMagnet(magnetUri string) (*Torrent, error) {
	res, err := ut.btClient.Parse(context.Background(), &prpc.DownloadRequest{
		Type:    prpc.DownloadRequest_MagnetUri,
		Content: []byte(magnetUri),
	})
	if err != nil || len(res.InfoHash.Hash) == 0 {
		return nil, errors.New("failed to parse magnet uri")
	}

	infoHash := TranInfoHash(res.InfoHash)

	ut.mtx.Lock()
	defer ut.mtx.Unlock()

	if t, ok := ut.torrents[*infoHash]; ok {
		return t, errors.New("existed torrent")
	}
	t, err := loadTorrentByInfoHash(&ut.btClient, infoHash)
	if err == nil {
		if len(t.GetMagnetUri()) == 0 {
			t.UpdateMagnetUri(magnetUri)
		}
		return t, errors.New("existed torrent")
	}
	t, err = newTorrent(&ut.btClient, infoHash, magnetUri)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (ut *UserTorrentsImpl) initTorrentLocked(infoHash *InfoHash, magnetUri string) *Torrent {
	if t, ok := ut.torrents[*infoHash]; ok {
		return t
	}

	t, err := loadTorrentByInfoHash(&ut.btClient, infoHash)
	if err != nil {
		t, err = newTorrent(&ut.btClient, infoHash, magnetUri)
		if err != nil {
			return nil
		}
	}
	if t == nil {
		return nil
	}
	ut.torrents[t.base.InfoHash] = t
	return t
}

func (ut *UserTorrentsImpl) HasTorrent(userId ptype.UserID, infoHash *InfoHash) bool {
	if userId == AdminId {
		return true
	}
	ut.mtx.Lock()
	t, ok1 := ut.torrents[*infoHash]
	u, ok2 := ut.users[userId]
	ut.mtx.Unlock()
	if !ok1 || !ok2 {
		return false
	}
	return u.hasTorrent(t.base.Id)
}

func (ut *UserTorrentsImpl) getTorrentLocked(infoHash *InfoHash) (*Torrent, error) {
	t, ok := ut.torrents[*infoHash]
	if !ok {
		return nil, errors.New("not found bt")
	}
	return t, nil
}

func (ut *UserTorrentsImpl) GetTorrent(infoHash *InfoHash) (*Torrent, error) {
	ut.mtx.Lock()
	defer ut.mtx.Unlock()
	return ut.getTorrentLocked(infoHash)
}

func (ut *UserTorrentsImpl) Download(params *DownloadParams) (*prpc.DownloadRespone, error) {
	req := params.Req
	res, err := ut.btClient.Parse(context.Background(), req)
	if err != nil {
		return nil, err
	}
	if len(res.InfoHash.Hash) == 0 {
		return nil, errors.New("invalid torrent")
	}

	infoHash := TranInfoHash(res.InfoHash)

	var magnetUri string
	if params.Req.Type == prpc.DownloadRequest_MagnetUri {
		magnetUri = string(params.Req.Content)
	}

	ut.mtx.Lock()
	defer ut.mtx.Unlock()

	res, err = ut.download(infoHash, req)
	if err == nil {
		log.Infof("[bt] uid:%d downloading %s", params.UserId, hex.EncodeToString([]byte(infoHash.Hash)))
		t := ut.initTorrentLocked(infoHash, magnetUri)
		if t != nil {
			t.addUser(ut.getUserDataLocked(params.UserId))
		}
		if params.UserId != ptype.AdminId {
			t.addUser(ut.getUserDataLocked(ptype.AdminId))
		}
		return res, nil
	} else {
		log.Infof("[bt] uid:%d failed to download %s err: %v", params.UserId, hex.EncodeToString([]byte(infoHash.Hash)), err)
		return res, status.Error(codes.InvalidArgument, "")
	}
}

func (ut *UserTorrentsImpl) NewDownloadTask(params *DownloadTaskParams) (*prpc.DownloadRespone, error) {
	req := params.Req
	res, err := ut.btClient.Parse(context.Background(), req)
	if err != nil {
		return nil, err
	}
	if len(res.InfoHash.Hash) == 0 {
		return nil, errors.New("invalid torrent")
	}

	infoHash := TranInfoHash(res.InfoHash)

	var magnetUri string
	if params.Req.Type == prpc.DownloadRequest_MagnetUri {
		magnetUri = string(params.Req.Content)
	}

	ut.mtx.Lock()
	defer ut.mtx.Unlock()

	res, err = ut.download(infoHash, req)
	if err == nil {
		log.Infof("[bt] uid:%d downloading %s", params.UserId, hex.EncodeToString([]byte(infoHash.Hash)))
		t := ut.initTorrentLocked(infoHash, magnetUri)
		if t != nil {
			t.addTask(ut.getUserDataLocked(params.UserId), params)
		}
		if params.UserId != ptype.AdminId {
			t.addUser(ut.getUserDataLocked(ptype.AdminId))
		}
		return res, nil
	} else {
		log.Infof("[bt] uid:%d failed to download %s err: %v", params.UserId, hex.EncodeToString([]byte(infoHash.Hash)), err)
		return res, status.Error(codes.InvalidArgument, "")
	}
}

func (ut *UserTorrentsImpl) RemoveTorrent(params *RemoveTorrentParams) (*prpc.RemoveTorrentRes, error) {
	infoHash := *TranInfoHash(params.Req.InfoHash)
	ut.mtx.Lock()
	defer ut.mtx.Unlock()
	if params.UserId == ptype.AdminId {
		log.Infof("[bt] uid:%d removing %s", params.UserId, hex.EncodeToString([]byte(infoHash.Hash)))
		res, err := ut.btClient.RemoveTorrent(context.Background(), params.Req)
		if err != nil {
			return nil, err
		}
		t, ok := ut.torrents[infoHash]
		if !ok {
			return nil, errors.New("not found torrent")
		}
		delete(ut.torrents, infoHash)
		t.remove()
		return res, err
	} else {
		t, ok := ut.torrents[infoHash]
		if !ok {
			return nil, errors.New("not found torrent")
		}
		t.removeUser(params.UserId)
	}
	return &prpc.RemoveTorrentRes{}, nil
}

func (ut *UserTorrentsImpl) GetMagnetUri(params *GetMagnetUriParams) (*prpc.GetMagnetUriRsp, error) {
	return ut.btClient.GetMagnetUri(context.Background(), params.Req)
}

func (ut *UserTorrentsImpl) GetTorrents(userId ptype.UserID) []*Torrent {
	if userId == ptype.AdminId {
		ut.mtx.Lock()
		defer ut.mtx.Unlock()
		ret := []*Torrent{}
		for _, t := range ut.torrents {
			ret = append(ret, t)
		}
		return ret
	}
	return ut.getUserData(userId).getTorrents()
}

func (ut *UserTorrentsImpl) GetPeerInfo(params *GetPeerInfoParams) (*prpc.GetPeerInfoRsp, error) {
	if params.Req.InfoHash == nil && params.UserId != ptype.AdminId {
		return nil, errors.New("null infohash")
	}
	if !ut.HasTorrent(params.UserId, TranInfoHash(params.Req.InfoHash)) {
		return nil, errors.New("no auth")
	}
	return ut.btClient.GetPeerInfo(context.Background(), params.Req)
}
