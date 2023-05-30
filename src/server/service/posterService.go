package service

import (
	"bufio"
	"net/http"
	"os"
	"pnas/category"
	"pnas/prpc"
	"pnas/service/session"
	"pnas/setting"
	"pnas/user"
	"strconv"

	"github.com/gorilla/mux"
)

type PosterService struct {
	um       *user.UserManger
	shares   SharesInterface
	sessions session.SessionsInterface
	router   *mux.Router
}
type NewPosterServiceParams struct {
	UserManger *user.UserManger
	Shares     SharesInterface
	Sessions   session.SessionsInterface
	Router     *mux.Router
}

func newPosterService(params *NewPosterServiceParams) *PosterService {
	ps := &PosterService{
		um:       params.UserManger,
		shares:   params.Shares,
		sessions: params.Sessions,
		router:   params.Router,
	}
	ps.registerUrl()
	return ps
}

func (v *PosterService) registerUrl() {
	v.router.Handle("/item", http.HandlerFunc(v.handlerItemPoster))
}

func (p *PosterService) handlerItemPoster(w http.ResponseWriter, r *http.Request) {
	s, _ := p.sessions.GetSession(r)
	queryParams := r.URL.Query()
	var userId user.ID
	itemIdTmp, _ := strconv.ParseInt(queryParams.Get("itemid"), 10, 64)
	itemId := category.ID(itemIdTmp)
	if s == nil {
		var err error
		userId, itemId, err = GetSharedItemInfo(p.shares, p.um, queryParams)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		userId = s.UserId
	}
	item, err := p.um.QueryItem(userId, itemId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pp := item.GetPosterPath()
	itype := item.GetType()
	if itype == prpc.CategoryItem_Directory {
		pp = setting.GS.Server.PosterPath + "/default_folder.png"
	} else if itype == prpc.CategoryItem_Home {
		pp = setting.GS.Server.PosterPath + "/house.png"
	} else {
		pp = setting.GS.Server.PosterPath + "/" + pp
	}

	if len(pp) > 0 {
		f, err := os.Open(pp)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer f.Close()
		bufio.NewReader(f).WriteTo(w)
	}
}
