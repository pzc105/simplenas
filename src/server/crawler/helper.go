package crawler

import (
	"pnas/ptype"
	"pnas/user"
	"sync/atomic"
)

func fetchCategoryId(name string, magnetShares user.IMagnetSharesService) (ptype.CategoryID, error) {
	items, _ := magnetShares.QueryMagnetCategorys(&user.QueryCategoryParams{
		ParentId:     magnetShares.GetMagnetRootId(),
		CategoryName: name,
	})

	var rid ptype.CategoryID
	var stop atomic.Bool
	stop.Store(false)

	if len(items) == 0 {
		var err error
		rid, err = magnetShares.AddMagnetCategory(&user.AddMagnetCategoryParams{
			ParentId:  magnetShares.GetMagnetRootId(),
			Name:      name,
			Introduce: "from crawler",
			Creator:   ptype.AdminId,
		})
		if err != nil {
			return -1, err
		}
	} else {
		rid = items[0].GetItemBaseInfo().Id
	}
	return rid, nil
}