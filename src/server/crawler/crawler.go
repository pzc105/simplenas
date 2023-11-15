package crawler

import (
	"pnas/bt"
	"pnas/log"
	"pnas/ptype"
	"pnas/user"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gocolly/colly"
)

const (
	CategoryName = "36dm"
)

func Go36dmBackgroup(magnetShares user.IMagnetSharesService, ut bt.UserTorrents) {
	items, _ := magnetShares.QueryMagnetCategorys(&user.QueryCategoryParams{
		ParentId:     magnetShares.GetMagnetRootId(),
		CategoryName: CategoryName,
	})

	var rid ptype.CategoryID
	var stop atomic.Bool
	stop.Store(false)

	if len(items) == 0 {
		var err error
		rid, err = magnetShares.AddMagnetCategory(&user.AddMagnetCategoryParams{
			ParentId:  magnetShares.GetMagnetRootId(),
			Name:      CategoryName,
			Introduce: "from crawler",
			Creator:   ptype.AdminId,
		})
		if err != nil {
			log.Errorf("failed to create crawler category: %v", err)
		}
	} else {
		rid = items[0].GetItemBaseInfo().Id
	}

	c := colly.NewCollector(
		colly.Async(true),
		colly.URLFilters(
			regexp.MustCompile(`https://www\.36dm\.org/forum-1.*`),
			regexp.MustCompile(`https://www\.36dm\.org/thread.*`),
		),
	)

	var flagMtx sync.Mutex
	flag := make(map[string]bool)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		var Name string
		var Uri string
		e.ForEach("div.media-body h4", func(i int, e *colly.HTMLElement) {
			Name = strings.Trim(e.Text, " \t\n")
		})
		e.ForEach("a[href]", func(i int, e *colly.HTMLElement) {
			link := e.Attr("href")
			if len(link) == 0 {
				return
			}
			if strings.HasPrefix(link, "magnet:?") {
				Uri = link
				return
			} else if strings.Contains(e.Request.URL.Path, "thread") {
				return
			}

			flagMtx.Lock()
			if _, ok := flag[link]; ok {
				flagMtx.Unlock()
				return
			}
			flag[link] = true
			flagMtx.Unlock()
			c.Visit(e.Request.AbsoluteURL(link))
		})
		if len(Uri) == 0 {
			return
		}

		t, err := ut.NewTorrentByMagnet(Uri)

		if err == nil {
			magnetShares.AddMagnetUri(&user.AddMagnetUriParams{
				T:          t,
				CategoryId: rid,
				Name:       "",
				Introduce:  Name,
				Creator:    ptype.AdminId,
			})
		}
	})

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})

	c.Visit("https://www.36dm.org/forum-1.htm")

	c.Wait()

	log.Info("Go36dmBackgroup done. restart...")

	timer := time.NewTimer(time.Hour * 3)
	<-timer.C
	go Go36dmBackgroup(magnetShares, ut)
}
