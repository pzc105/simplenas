package bt

import (
	"context"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/setting"
	"sync"

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
	mtx       sync.Mutex
	userId    ptype.UserID
	torrents  map[ptype.TorrentID]bool
	callbacks map[ptype.SessionID]UserOnBtStatusCallback
}

func (ud *userData) init() {
	ud.torrents = make(map[ptype.TorrentID]bool)
	ud.callbacks = make(map[ptype.SessionID]UserOnBtStatusCallback)
}

func (ud *userData) setCallback(sid ptype.SessionID, callback UserOnBtStatusCallback) {
	ud.mtx.Lock()
	if callback == nil {
		delete(ud.callbacks, sid)
	} else {
		ud.callbacks[sid] = callback
	}
	ud.mtx.Unlock()
}

func (ud *userData) getCallbackLocked() []UserOnBtStatusCallback {
	var ret []UserOnBtStatusCallback
	for _, v := range ud.callbacks {
		ret = append(ret, v)
	}
	return ret
}

func (ud *userData) onBtStatus(id ptype.TorrentID, s *prpc.TorrentStatus) {
	if !ud.hasTorrent(id) {
		return
	}
	ud.mtx.Lock()
	callbacks := ud.getCallbackLocked()
	ud.mtx.Unlock()
	for _, callback := range callbacks {
		if callback != nil {
			callback(s)
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

func (ud *userData) initTorrent(id ptype.TorrentID) {
	ud.mtx.Lock()
	ud.torrents[id] = true
	ud.mtx.Unlock()
}

func (ud *userData) addTorrent(id ptype.TorrentID) error {
	if ud.hasTorrent(id) {
		return errors.New(("duplicated"))
	}
	ud.mtx.Lock()
	defer ud.mtx.Unlock()
	sql := "insert into user_torrent (user_id, torrent_id) values(?, ?)"
	_, err := db.Exec(sql, ud.userId, id)
	if err == nil {
		ud.torrents[id] = true
	}
	return err
}

func (ud *userData) removeTorrent(id ptype.TorrentID, dodb bool) error {
	ud.mtx.Lock()
	defer ud.mtx.Unlock()
	_, ok := ud.torrents[id]
	if !ok {
		return errors.New("not found torrent")
	}
	if dodb {
		sql := "delete from user_torrent where user_id=? and torrent_id=?"
		_, err := db.Exec(sql, ud.userId, id)
		return err
	}
	return nil
}

type UserTorrentsImpl struct {
	mtx      sync.Mutex
	torrents map[InfoHash]*Torrent
	users    map[ptype.UserID]*userData

	btClient BtClient
}

func (ut *UserTorrentsImpl) Init() {
	ut.torrents = make(map[InfoHash]*Torrent)
	ut.users = make(map[ptype.UserID]*userData)
	ut.load()

	ut.btClient.Init(WithOnStatus(ut.onBtStatus),
		WithOnConnect(ut.handleBtClientConnected),
		WithOnFileCompleted(ut.handleBtFileCompleted))
}

func (ut *UserTorrentsImpl) load() {
	sql := `select user_id, torrent_id from user_torrent`
	rows, err := db.Query(sql)
	if err != nil {
		log.Warnf("load user torrent err: %v", err)
		return
	}
	defer rows.Close()
	flag := make(map[ptype.TorrentID]bool)
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
		if _, ok = flag[tid]; !ok {
			flag[tid] = true
			t = loadTorrent(&ut.btClient, tid)
			if t == nil {
				ut.mtx.Unlock()
				continue
			}
			ut.torrents[t.base.InfoHash] = t
			ut.mtx.Unlock()
			t.addUser(uid)
			ut.initUserTorrent(t, uid)
		} else {
			ut.mtx.Unlock()
		}
	}
}

func (ut *UserTorrentsImpl) Close() {
	ut.btClient.Close()
}

func (ut *UserTorrentsImpl) handleBtClientConnected() {
	log.Info("connected to bt service")
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
		req.SavePath = setting.GS().Bt.SavePath
		_, err = ut.btClient.Download(context.Background(), req)
		if err != nil {
			log.Warnf("failed download err: %v", err)
		}
	}
}

func (ut *UserTorrentsImpl) SetCallback(userId ptype.UserID, sid ptype.SessionID, callback UserOnBtStatusCallback) {
	ut.getUserData(userId).setCallback(sid, callback)
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
	// if !t.hasBaseInfo() {
	// 	req := &prpc.GetTorrentInfoReq{
	// 		InfoHash: s.InfoHash,
	// 	}
	// 	res, err := ut.btClient.GetTorrentInfo(context.Background(), req)
	// 	if err == nil {
	// 		t.updateTorrentInfo(res.TorrentInfo)
	// 	}
	// }
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

func (ut *UserTorrentsImpl) getUserData(uid ptype.UserID) *userData {
	ut.mtx.Lock()
	defer ut.mtx.Unlock()
	ud, ok := ut.users[uid]
	if ok {
		return ud
	}
	ud = &userData{
		userId: uid,
	}
	ud.init()
	ut.users[uid] = ud
	return ud
}

func (ut *UserTorrentsImpl) initTorrent(infoHash *InfoHash) *Torrent {
	ut.mtx.Lock()
	defer ut.mtx.Unlock()

	if t, ok := ut.torrents[*infoHash]; ok {
		return t
	}

	// TODO: handle mysql error
	t := loadTorrentByInfoHash(&ut.btClient, infoHash)
	if t == nil {
		t = newTorrent(&ut.btClient, infoHash)
	}
	if t == nil {
		return nil
	}
	ut.torrents[t.base.InfoHash] = t
	return t
}

func (ut *UserTorrentsImpl) initUserTorrent(t *Torrent, uid ptype.UserID) {
	ut.getUserData(uid).initTorrent(t.base.Id)
}

func (ut *UserTorrentsImpl) saveUserTorrent(t *Torrent, uid ptype.UserID) error {
	ud := ut.getUserData(uid)
	ud.addTorrent(t.base.Id)
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

func (ut *UserTorrentsImpl) GetTorrent(infoHash *InfoHash) (*Torrent, error) {
	ut.mtx.Lock()
	defer ut.mtx.Unlock()
	t, ok := ut.torrents[*infoHash]
	if !ok {
		return nil, errors.New("not found bt")
	}
	return t, nil
}

func (ut *UserTorrentsImpl) Download(params *DownloadParams) (*prpc.DownloadRespone, error) {
	req := params.Req
	res, err := ut.btClient.Parse(context.Background(), req)
	if err != nil {
		return nil, err
	}

	infoHash := TranInfoHash(res.InfoHash)

	t, err := ut.GetTorrent(infoHash)
	ud := ut.getUserData(params.UserId)
	if err == nil {
		if ud.hasTorrent(t.base.Id) {
			return nil, errors.New("duplicated")
		}
	}

	resumeData, err := loadResumeData(infoHash)
	if err == nil {
		req.Type = prpc.DownloadRequest_Resume
		req.Content = resumeData
	}

	req.SavePath = setting.GS().Bt.SavePath
	res, err = ut.btClient.Download(context.Background(), req)
	if err == nil {
		t := ut.initTorrent(infoHash)
		if t != nil {
			t.addUser(params.UserId)
			ut.saveUserTorrent(t, params.UserId)
		}
		return res, nil
	} else {
		return res, status.Error(codes.InvalidArgument, "")
	}
}

func (ut *UserTorrentsImpl) RemoveTorrent(params *RemoveTorrentParams) (*prpc.RemoveTorrentRes, error) {
	if params.UserId == ptype.AdminId {
		res, err := ut.btClient.RemoveTorrent(context.Background(), params.Req)
		if err != nil {
			return nil, err
		}
		ut.mtx.Lock()
		t, ok := ut.torrents[*TranInfoHash(params.Req.InfoHash)]
		if !ok {
			ut.mtx.Unlock()
			return nil, errors.New("not found torrent")
		}
		uids := t.getAllUser()
		for _, uid := range uids {
			ud := ut.users[uid]
			t.removeUser(uid)
			ud.removeTorrent(t.base.Id, false)
		}
		err = deleteUserTorrentids(t.base.Id, uids...)
		ut.mtx.Unlock()
		return res, err
	} else {
		ut.mtx.Lock()
		t, ok := ut.torrents[*TranInfoHash(params.Req.InfoHash)]
		if !ok {
			ut.mtx.Unlock()
			return nil, errors.New("not found torrent")
		}
		ut.mtx.Unlock()
		ud := ut.getUserData(params.UserId)
		t.removeUser(params.UserId)
		ud.removeTorrent(t.base.Id, true)
	}
	return &prpc.RemoveTorrentRes{}, nil
}

func (ut *UserTorrentsImpl) GetMagnetUri(params *GetMagnetUriParams) (*prpc.GetMagnetUriRsp, error) {
	return ut.btClient.GetMagnetUri(context.Background(), params.Req)
}
