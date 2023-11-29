package bt

import (
	"errors"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/utils"
	"sync"
)

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
