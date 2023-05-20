package service

import (
	"bufio"
	"net/http"
	"os"
	"pnas/category"
	"pnas/prpc"
	"pnas/setting"
	"strconv"

	"github.com/gorilla/mux"
)

type PosterService struct {
	coreSer CoreServiceInterface
	router  *mux.Router
}

func newPosterService(core CoreServiceInterface, router *mux.Router) *PosterService {
	ps := &PosterService{
		coreSer: core,
	}
	ps.router = router
	ps.registerUrl()
	return ps
}

func (v *PosterService) registerUrl() {
	v.router.Handle("/item/{itemId}", http.HandlerFunc(v.handlerItemPoster))
}

func (p *PosterService) handlerItemPoster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemIdStr, ok := vars["itemId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s := p.coreSer.GetSession(r)
	if s == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	item, err := p.coreSer.GetUserManager().QueryItem(s.UserId, category.ID(itemId))
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
