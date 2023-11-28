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

type userData struct {
	mtx        sync.Mutex
	userId     ptype.UserID
	torrents   map[ptype.TorrentID]*Torrent
	callbacks  map[ptype.SessionID]UserOnBtStatusCallback
	callbacks2 map[ptype.TorrentID]map[ptype.TaskId]UserOnBtStatusCallback

	callbackTaskQueue *utils.TaskQueue
}

func (ud *userData) init(taskQueue *utils.TaskQueue) {
	ud.torrents = make(map[ptype.TorrentID]*Torrent)
	ud.callbacks = make(map[ptype.SessionID]UserOnBtStatusCallback)
	ud.callbacks2 = make(map[ptype.TorrentID]map[ptype.TaskId]UserOnBtStatusCallback)
	ud.callbackTaskQueue = taskQueue
}

func (ud *userData) setTaskCallback(params *SetTaskCallbackParams) {
	ud.mtx.Lock()
	defer ud.mtx.Unlock()
	cbs, ok := ud.callbacks2[params.TorrentId]
	if !ok {
		return
	}
	if params.Callback == nil {
		delete(cbs, params.TaskId)
	} else {
		cbs[params.TaskId] = params.Callback
	}
}

func (ud *userData) setSessionCallback(sid ptype.SessionID, callback UserOnBtStatusCallback) {
	ud.mtx.Lock()
	defer ud.mtx.Unlock()
	if callback == nil {
		delete(ud.callbacks, sid)
	} else {
		ud.callbacks[sid] = callback
	}
}

func (ud *userData) getCallbackLocked(tid ptype.TorrentID) []UserOnBtStatusCallback {
	var ret []UserOnBtStatusCallback
	for _, v := range ud.callbacks {
		ret = append(ret, v)
	}
	for _, v := range ud.callbacks2[tid] {
		ret = append(ret, v)
	}
	return ret
}

func (ud *userData) onBtStatus(tid ptype.TorrentID, s *prpc.TorrentStatus) {
	if !ud.hasTorrent(tid) {
		return
	}
	ud.mtx.Lock()
	callbacks := ud.getCallbackLocked(tid)
	ud.mtx.Unlock()
	for _, callback := range callbacks {
		c := callback
		if c != nil {
			ud.callbackTaskQueue.Put(func() {
				c(nil, s)
			})
		}
	}
}

func (ud *userData) hasTorrent(id ptype.TorrentID) bool {
	if ud.userId == ptype.AdminId {
		return true
	}
	ud.mtx.Lock()
	_, ok := ud.torrents[id]
	ud.mtx.Unlock()
	return ok
}

func (ud *userData) initTorrentLocked(t *Torrent) {
	ud.torrents[t.base.Id] = t
	ud.callbacks2[t.base.Id] = make(map[ptype.TaskId]UserOnBtStatusCallback)
}

func (ud *userData) initTorrent(t *Torrent) {
	ud.mtx.Lock()
	ud.initTorrentLocked(t)
	ud.mtx.Unlock()
}

func (ud *userData) addTorrent(t *Torrent) error {
	ud.mtx.Lock()
	defer ud.mtx.Unlock()
	sql := "insert into user_torrent (user_id, torrent_id) values(?, ?)"
	_, err := db.Exec(sql, ud.userId, t.base.Id)
	if err == nil {
		ud.initTorrentLocked(t)
	}
	return err
}

func (ud *userData) removeTorrent(id ptype.TorrentID, dodb bool) error {
	taskCallbacks := []UserOnBtStatusCallback{}
	ud.mtx.Lock()
	_, ok := ud.torrents[id]
	if !ok {
		ud.mtx.Unlock()
		return errors.New("not found torrent")
	}

	for _, c := range ud.callbacks2[id] {
		taskCallbacks = append(taskCallbacks, c)
	}
	delete(ud.callbacks2, id)
	delete(ud.torrents, id)
	if dodb {
		sql := "delete from user_torrent where user_id=? and torrent_id=?"
		_, err := db.Exec(sql, ud.userId, id)
		if err != nil {
			log.Warnf("[bt] failed delete user_torrent er:%v", err)
		}
	}
	ud.mtx.Unlock()
	err := errors.New("removing torrent")
	for _, callback := range taskCallbacks {
		c := callback
		if c != nil {
			ud.callbackTaskQueue.Put(func() {
				c(err, nil)
			})
		}
	}
	return nil
}

func (ud *userData) getTorrents() []*Torrent {
	ud.mtx.Lock()
	defer ud.mtx.Unlock()
	ts := []*Torrent{}
	for _, t := range ud.torrents {
		ts = append(ts, t)
	}
	return ts
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
	sql := `select user_id, torrent_id from user_torrent`
	rows, err := db.Query(sql)
	if err != nil {
		log.Warnf("load user torrent err: %v", err)
		return
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
		ut.mtx.Unlock()
		t.addUser(uid)
		ut.initUserTorrent(t, uid)
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
	uids := t.getAllUser()
	for _, uid := range uids {
		ut.getUserData(uid).onBtStatus(t.base.Id, s)
	}
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

func (ut *UserTorrentsImpl) initUserTorrent(t *Torrent, uid ptype.UserID) {
	ut.getUserData(uid).initTorrent(t)
}

func (ut *UserTorrentsImpl) saveUserTorrentLocked(t *Torrent, uid ptype.UserID) error {
	ud := ut.getUserDataLocked(uid)
	ud.addTorrent(t)
	return nil
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
	t, err := ut.getTorrentLocked(infoHash)
	if err == nil {
		if t.GetState() == prpc.BtStateEnum_seeding {
			log.Debugf("[bt] user %d add %s", params.UserId, hex.EncodeToString([]byte(infoHash.Hash)))
			t.addUser(params.UserId)
			ut.saveUserTorrentLocked(t, params.UserId)
			return &prpc.DownloadRespone{
				InfoHash: res.InfoHash,
			}, ErrDownloaded
		}
	}

	res, err = ut.download(infoHash, req)
	if err == nil {
		log.Debugf("[bt] user %d downloading %s", params.UserId, hex.EncodeToString([]byte(infoHash.Hash)))
		t := ut.initTorrentLocked(infoHash, magnetUri)
		if t != nil {
			t.addUser(params.UserId)
			ut.saveUserTorrentLocked(t, params.UserId)
		}

		if params.UserId != ptype.AdminId {
			t.addUser(ptype.AdminId)
			ut.saveUserTorrentLocked(t, ptype.AdminId)
		}
		return res, nil
	} else {
		log.Infof("[bt] user %d failed to download %s err: %v", params.UserId, hex.EncodeToString([]byte(infoHash.Hash)), err)
		return res, status.Error(codes.InvalidArgument, "")
	}
}

func (ut *UserTorrentsImpl) RemoveTorrent(params *RemoveTorrentParams) (*prpc.RemoveTorrentRes, error) {
	infoHash := *TranInfoHash(params.Req.InfoHash)
	log.Infof("[bt] user %d removing %s", params.UserId, hex.EncodeToString([]byte(infoHash.Hash)))
	ut.mtx.Lock()
	defer ut.mtx.Unlock()
	if params.UserId == ptype.AdminId {
		res, err := ut.btClient.RemoveTorrent(context.Background(), params.Req)
		if err != nil {
			return nil, err
		}
		t, ok := ut.torrents[infoHash]
		if !ok {
			return nil, errors.New("not found torrent")
		}
		delete(ut.torrents, infoHash)
		// TODO: fix concurrency issues
		t.delResumeData()
		uids := t.getAllUser()
		err = deleteUserTorrentids(t.base.Id, uids...)
		for _, uid := range uids {
			ud := ut.users[uid]
			t.removeUser(uid)
			ud.removeTorrent(t.base.Id, false)
		}
		return res, err
	} else {
		t, ok := ut.torrents[infoHash]
		if !ok {
			return nil, errors.New("not found torrent")
		}
		ud := ut.getUserDataLocked(params.UserId)
		t.removeUser(params.UserId)
		ud.removeTorrent(t.base.Id, true)
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
