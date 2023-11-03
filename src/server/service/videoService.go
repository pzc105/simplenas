package service

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"pnas/db"
	"pnas/log"
	"pnas/service/session"
	"pnas/setting"
	"pnas/user"
	"pnas/video"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/grafov/m3u8"
)

type UserVideoData interface {
	HasVideo(userId user.ID, vid video.ID) bool
}

type VideoService struct {
	ud       UserVideoData
	shares   IItemShares
	sessions session.ISessions
	router   *mux.Router
}

func saveStartTime(userId user.ID, vid video.ID, lastTime string) {
	db.GREDIS.Set(context.Background(), fmt.Sprintf("video_offset_%d_%d", userId, vid), lastTime, 0)
}

func loadStartTime(userId user.ID, vid video.ID) string {
	startTime, err := db.GREDIS.Get(context.Background(), fmt.Sprintf("video_offset_%d_%d", userId, vid)).Result()
	if err == nil {
		return startTime
	}
	return "0"
}

type NewVideoServiceParams struct {
	UserData UserVideoData
	Shares   IItemShares
	Sessions session.ISessions
	Router   *mux.Router
}

func newVideoService(params *NewVideoServiceParams) *VideoService {
	vs := &VideoService{
		ud:       params.UserData,
		shares:   params.Shares,
		sessions: params.Sessions,
		router:   params.Router,
	}
	vs.registerUrl()
	return vs
}

func (v *VideoService) getAccessUser(r *http.Request) (user.ID, error) {
	queryParams := r.URL.Query()
	shareid := queryParams.Get("shareid")
	if len(shareid) > 0 {
		si, err := v.shares.GetShareItemInfo(shareid)
		if err != nil {
			return -1, err
		}
		return si.UserId, nil
	} else {
		s, err := v.sessions.GetSession(r)
		if err != nil {
			return -1, err
		}
		return s.UserId, nil
	}
}

func (v *VideoService) checkAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		vidstr, ok := vars["vid"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		vid, err := strconv.ParseInt(vidstr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userId, err := v.getAccessUser(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !v.ud.HasVideo(userId, video.ID(vid)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (v *VideoService) registerUrl() {
	v.router.Use(v.checkAuth)
	v.router.Handle("/{vid}", http.HandlerFunc(v.handlerHlsMasterList))
	v.router.Handle("/{vid}/subtitle/{caption}", http.HandlerFunc(v.handlerSubtitle))
	v.router.Handle("/{vid}/poster", http.HandlerFunc(v.handlePoster))
	v.router.Handle("/{vid}/set_offsettime/{offset}", http.HandlerFunc(v.handleSetOffsetTime))
	v.router.Handle("/{vid}/get_offsettime", http.HandlerFunc(v.handleGetOffsetTime))
	v.router.Handle("/{vid}/stream_{sid}/segment/{segment}", http.HandlerFunc(v.handlerHlsSegment))
	v.router.Handle("/{vid}/stream_{sid}/{playlist}", http.HandlerFunc(v.handlerHlsPlayList))
}

func (v *VideoService) handlerHlsMasterList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vid, ok := vars["vid"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	playlistPath := setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%s/master.m3u8", vid)
	f, err := os.Open(playlistPath)
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	masterpl := m3u8.NewMasterPlaylist()
	err = masterpl.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	queryParams := r.URL.Query()
	itemid := queryParams.Get("itemid")
	shareid := queryParams.Get("shareid")
	for _, v := range masterpl.Variants {
		v.URI = fmt.Sprintf("/video/%s/%s", vid, v.URI)
		if len(shareid) > 0 && len(itemid) > 0 {
			v.URI += fmt.Sprintf("?shareid=%s&itemid=%s", shareid, itemid)
		}
		for _, al := range v.Alternatives {
			if al == nil {
				continue
			}
			al.URI = fmt.Sprintf("/video/%s/%s", vid, al.URI)
			if len(shareid) > 0 && len(itemid) > 0 {
				al.URI += fmt.Sprintf("?shareid=%s&itemid=%s", shareid, itemid)
			}
		}
	}
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	masterpl.Encode().WriteTo(w)
}

func (v *VideoService) handlerHlsPlayList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vid, ok1 := vars["vid"]
	sid, ok2 := vars["sid"]
	pn, ok3 := vars["playlist"]
	if !ok1 || !ok2 || !ok3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	playlistPath := setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%s/stream_%s/%s", vid, sid, pn)
	f, err := os.Open(playlistPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pl, err := m3u8.NewMediaPlaylist(0, 10)

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	err = pl.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if pl.Map != nil {
		pl.Map.URI = fmt.Sprintf("/video/%s/stream_%s/segment/%s", vid, sid, pl.Map.URI)
	}
	queryParams := r.URL.Query()
	itemid := queryParams.Get("itemid")
	shareid := queryParams.Get("shareid")
	for _, v := range pl.Segments {
		if v == nil {
			continue
		}
		v.URI = fmt.Sprintf("/video/%s/stream_%s/segment/%s", vid, sid, v.URI)
		if len(shareid) > 0 && len(itemid) > 0 {
			v.URI += fmt.Sprintf("?shareid=%s&itemid=%s", shareid, itemid)
		}
	}
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	pl.Encode().WriteTo(w)
}

func (v *VideoService) handlerHlsSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vid, ok1 := vars["vid"]
	sid, ok2 := vars["sid"]
	sg, ok3 := vars["segment"]
	if !ok1 || !ok2 || !ok3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	segmentPath := setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%s/stream_%s/%s", vid, sid, sg)
	f, err := os.Open(segmentPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "video/mp2t")
	bufio.NewReader(f).WriteTo(w)
}

func (v *VideoService) handlerSubtitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vid, ok1 := vars["vid"]
	cid, ok2 := vars["caption"]
	if !ok1 || !ok2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	captionPath := setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%s/%s", vid, cid)
	f, err := os.Open(captionPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer f.Close()
	bufio.NewReader(f).WriteTo(w)
}

func (v *VideoService) handlePoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vid, ok1 := vars["vid"]
	if !ok1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	posterPath := setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%s/poster.png", vid)
	f, err := os.Open(posterPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer f.Close()
	bufio.NewReader(f).WriteTo(w)
}

func (v *VideoService) handleSetOffsetTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vidstr, ok1 := vars["vid"]
	if !ok1 {
		return
	}
	vid, err := strconv.ParseInt(vidstr, 10, 64)
	if err != nil {
		return
	}
	s, _ := v.sessions.GetSession(r)
	if s == nil || !v.ud.HasVideo(s.UserId, video.ID(vid)) {
		return
	}
	timeoffset, ok2 := vars["offset"]
	if !ok2 {
		return
	}
	saveStartTime(s.UserId, video.ID(vid), timeoffset)
}

func (v *VideoService) handleGetOffsetTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vidstr, ok1 := vars["vid"]
	if !ok1 {
		return
	}
	vid, err := strconv.ParseInt(vidstr, 10, 64)
	if err != nil {
		return
	}
	s, _ := v.sessions.GetSession(r)
	if s == nil || !v.ud.HasVideo(s.UserId, video.ID(vid)) {
		return
	}
	w.Write([]byte(loadStartTime(s.UserId, video.ID(vid))))
}
