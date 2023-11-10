package crawler

import (
	"context"
	"pnas/bt"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/user"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gocolly/colly"
)

const (
	CategoryName = "36dm"
)

func Go36dmBackgroup(magnetShares user.IMagnetSharesService, btClient *bt.BtClient) {
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
		rid = items[0].GetItemInfo().Id
	}

	c := colly.NewCollector(
		colly.Async(true),
		colly.URLFilters(
			regexp.MustCompile(`https://www\.36dm\.org.*`),
		),
	)

	var flagMtx sync.Mutex
	flag := make(map[int]bool)

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
			pi := strings.Index(link, "thread-")
			if pi != -1 {
				ns := link[pi+len("thread-"):]
				ns, f := strings.CutSuffix(ns, ".htm")
				if f {
					num, err := strconv.Atoi(ns)
					if err == nil {
						flagMtx.Lock()
						if _, ok := flag[num]; ok {
							flagMtx.Unlock()
							return
						}
						flag[num] = true
						flagMtx.Unlock()
						c.Visit(e.Request.AbsoluteURL(link))
					}
				}
				return
			}
			c.Visit(e.Request.AbsoluteURL(link))
		})
		if len(Uri) == 0 {
			return
		}

		rsp, err := btClient.Parse(context.Background(), &prpc.DownloadRequest{
			Type:    prpc.DownloadRequest_MagnetUri,
			Content: []byte(Uri),
		})

		if err != nil {
			return
		}

		sql := "insert into magnet(version, info_hash, magnet_uri) values(?, ?, ?)"
		dr, err := db.Exec(sql, rsp.InfoHash.Version, rsp.InfoHash.Hash, Uri)
		if err != nil {
			return
		}
		af, err := dr.RowsAffected()
		if err != nil || af == 0 {
			return
		}

		magnetShares.AddMagnetUri(&user.AddMagnetUriParams{
			Uri:        Uri,
			CategoryId: rid,
			Name:       "",
			Introduce:  Name,
			Creator:    ptype.AdminId,
		})
	})

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 3})

	c.Visit("https://www.36dm.org")

	c.Wait()
}
