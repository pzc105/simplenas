package crawler

import (
	"pnas/category"
	"pnas/log"
	"pnas/user"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

const (
	CategoryName = "36dm"
)

func Go36dmBackgroup(magnetShares user.IMagnetSharesService) {

	items, _ := magnetShares.QueryMagnetCategorys(&user.QueryCategoryParams{
		ParentId:     magnetShares.GetMagnetRootId(),
		CategoryName: CategoryName,
	})

	var rid category.ID

	if len(items) == 0 {
		var err error
		rid, err = magnetShares.AddMagnetCategory(&user.AddMagnetCategoryParams{
			ParentId:  magnetShares.GetMagnetRootId(),
			Name:      CategoryName,
			Introduce: "from crawler",
			Creator:   category.AdminId,
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

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		var Name string
		var Uri string
		e.ForEach("div.media-body h4", func(i int, e *colly.HTMLElement) {
			Name = strings.Trim(e.Text, " \t\n")
		})
		e.ForEach("a[href]", func(i int, e *colly.HTMLElement) {
			link := e.Attr("href")
			if strings.HasPrefix(link, "magnet:?") {
				Uri = link
			}
		})
		if len(Uri) == 0 {
			return
		}
		magnetShares.AddMagnetUri(&user.AddMagnetUriParams{
			Uri:        Uri,
			CategoryId: rid,
			Name:       "",
			Introduce:  Name,
			Creator:    category.AdminId,
		})
	})

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})

	c.Visit("https://www.36dm.org")

	c.Wait()
}
