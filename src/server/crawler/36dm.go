package crawler

import (
	"pnas/log"
	"pnas/ptype"
	"pnas/user"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

const (
	categoryName36dm = "36dm"
)

func Go36dmBackgroup(magnetShares user.IMagnetSharesService, maxDepth int) {
	rid, err := fetchCategoryId(categoryName36dm, magnetShares)
	if err != nil {
		log.Warnf("can't fetch 36dm category err:%v", err)
		return
	}

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(maxDepth),
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
		e.ForEach("div.media-body h4", func(_ int, e *colly.HTMLElement) {
			Name = strings.Trim(e.Text, " \t\n")
		})
		e.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
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

		magnetShares.AddMagnetUri(&user.AddMagnetUriParams{
			CategoryId: rid,
			Name:       Name,
			Creator:    ptype.AdminId,
			Uri:        Uri,
		})
	})

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})

	c.Visit("https://www.36dm.org/forum-1.htm")

	c.Wait()

	log.Info("Go36dmBackgroup done. restart...")

	if maxDepth > 0 {
		timer := time.NewTimer(time.Hour * 3)
		<-timer.C
	}
	go Go36dmBackgroup(magnetShares, 2)
}
