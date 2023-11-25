package crawler

import (
	"pnas/log"
	"pnas/ptype"
	"pnas/user"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

const (
	categoryNameAcg = "acg"
)

type GoAcgBackgroupParams struct {
	MagnetShares user.IMagnetSharesService
	MaxDepth     int
	ProxyUrl     string
}

func GoAcgBackgroup(params *GoAcgBackgroupParams) {
	rid, err := fetchCategoryId(categoryNameAcg, params.MagnetShares)
	if err != nil {
		log.Warnf("can't fetch 36dm category err:%v", err)
		return
	}

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(params.MaxDepth),
		colly.URLFilters(
			regexp.MustCompile(`https://acg\.rip/1/page/.*`),
			regexp.MustCompile(`https://acg\.rip/t/.*`),
		),
	)
	if len(params.ProxyUrl) > 0 {
		err := c.SetProxy(params.ProxyUrl)
		if err != nil {
			log.Warnf("failed to set proxy err:%v", err)
		}
	}
	var flagMtx sync.Mutex
	flag := make(map[string]bool)
	names := make(map[int64]string)
	cacheTorrent := make(map[int64][]byte)

	const prefix = "/t/"
	const suffix = ".torrent"

	c.OnResponse(func(rsp *colly.Response) {
		if strings.HasPrefix(rsp.Request.URL.Path, prefix) && strings.HasSuffix(rsp.Request.URL.Path, suffix) {
			numstr := rsp.Request.URL.Path[len(prefix):]
			numstr = numstr[:len(numstr)-len(suffix)]
			num, err := strconv.ParseInt(numstr, 10, 64)
			if err == nil {
				flagMtx.Lock()
				name, ok := names[num]
				flagMtx.Unlock()
				if ok {
					params.MagnetShares.AddMagnetUriByTorrent(&user.AddMagnetUriParams{
						CategoryId: rid,
						Name:       name,
						Creator:    ptype.AdminId,
						Torrent:    rsp.Body,
					})
				} else {
					flagMtx.Lock()
					cacheTorrent[num] = rsp.Body
					flagMtx.Unlock()
				}
			}
		}
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Request.URL.Path, prefix) && !strings.HasSuffix(e.Request.URL.Path, suffix) {
			numstr := e.Request.URL.Path[len(prefix):]
			num, err := strconv.ParseInt(numstr, 10, 64)
			if err == nil {
				var Name string
				e.ForEach("div.panel-heading", func(_ int, e *colly.HTMLElement) {
					Name = strings.Trim(e.Text, " \t\n")
				})
				if len(Name) > 0 {
					flagMtx.Lock()
					torrentBytes, ok := cacheTorrent[num]
					flagMtx.Unlock()
					if ok {
						flagMtx.Lock()
						delete(cacheTorrent, num)
						flagMtx.Unlock()
						params.MagnetShares.AddMagnetUriByTorrent(&user.AddMagnetUriParams{
							CategoryId: rid,
							Name:       Name,
							Creator:    ptype.AdminId,
							Torrent:    torrentBytes,
						})
					} else {
						flagMtx.Lock()
						names[num] = Name
						flagMtx.Unlock()
					}
				}
			}
		}

		e.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
			link := e.Attr("href")
			if len(link) == 0 {
				return
			}
			absLink := e.Request.AbsoluteURL(link)
			flagMtx.Lock()
			if _, ok := flag[absLink]; ok {
				flagMtx.Unlock()
				return
			}
			flag[absLink] = true
			flagMtx.Unlock()
			c.Visit(absLink)
		})
	})

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 10})

	c.Visit("https://acg.rip/1/page/1")

	c.Wait()

	log.Info("GoAcgBackgroup done. restart...")

	if params.MaxDepth > 0 {
		timer := time.NewTimer(time.Hour * 3)
		<-timer.C
	}

	params.MaxDepth = 2
	go GoAcgBackgroup(params)
}
