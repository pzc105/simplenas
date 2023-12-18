package category

import (
	"fmt"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

const (
	RootId     = 1
	UsersId    = 3
	NotExisted = -2
)

type OtherInfo struct {
	MagnetUri string
}

type BaseItem struct {
	Id           ptype.CategoryID
	Creator      ptype.UserID
	CreatedAt    time.Time
	TypeId       prpc.CategoryItem_Type
	Name         string
	ResourcePath string
	PosterPath   string
	Introduce    string
	Other        OtherInfo
	UpdatedAt    time.Time
	ParentId     ptype.CategoryID
}

type CategoryItem struct {
	mtx        sync.Mutex
	base       BaseItem
	auth       utils.AuthBitSet
	subItemIds []ptype.CategoryID
}

func _initMagnet(item *CategoryItem) error {
	if item.base.TypeId == prpc.CategoryItem_MagnetUri {
		sql := `select magnet_uri from torrent where id=?`
		err := db.QueryRow(sql, item.base.ResourcePath).Scan(
			&item.base.Other.MagnetUri,
		)
		if err != nil {
			log.Warnf("[category] fialed to load magnet uri err: %v", err)
			return err
		}
		return nil
	}
	return errors.New("isn't magnet uri")
}

func _loadItem(itemId ptype.CategoryID) (*CategoryItem, error) {
	log.Debug("[category] loading item: ", itemId)
	var item CategoryItem
	item.base.Id = itemId
	var byteAuth []byte
	sql := `select type_id, name, creator, auth, resource_path, poster_path, introduce, created_at, updated_at, parent_id
				from pnas.category_items
				where id=?`
	err := db.QueryRow(sql, itemId).Scan(
		&item.base.TypeId, &item.base.Name, &item.base.Creator, &byteAuth, &item.base.ResourcePath, &item.base.PosterPath,
		&item.base.Introduce, &item.base.CreatedAt, &item.base.UpdatedAt, &item.base.ParentId,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	item.auth = utils.NewBitSet(AuthMax)
	item.auth.UnmarshalBinary(byteAuth)

	err = item._initSubItemIds()
	if err != nil {
		return nil, err
	}
	_initMagnet(&item)
	log.Debug("[category] loaded item: ", itemId)
	return &item, nil
}

func _loadItemIdByName(parentId ptype.CategoryID, name string) (ptype.CategoryID, error) {
	var ret ptype.CategoryID
	sql := `select id from pnas.category_items where parent_id=? and name=?`
	err := db.QueryRow(sql, parentId, name).Scan(&ret)
	return ret, err
}

func _loadItems(itemIds ...ptype.CategoryID) ([]*CategoryItem, error) {
	if len(itemIds) == 0 {
		return []*CategoryItem{}, nil
	}
	log.Debug("[category] loading items: ", itemIds)
	var conds []string
	for _, id := range itemIds {
		conds = append(conds, fmt.Sprintf("id=%d", id))
	}
	cond := strings.Join(conds, " or ")
	sql := `select id, type_id, name, creator, auth, resource_path, poster_path, introduce, created_at, updated_at, parent_id
					from pnas.category_items where ` + cond
	rows, err := db.Query(sql)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()
	var items []*CategoryItem
	var byteAuth []byte
	for rows.Next() {
		var item CategoryItem
		err = rows.Scan(
			&item.base.Id, &item.base.TypeId, &item.base.Name, &item.base.Creator, &byteAuth,
			&item.base.ResourcePath, &item.base.PosterPath,
			&item.base.Introduce, &item.base.CreatedAt, &item.base.UpdatedAt, &item.base.ParentId,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		item.auth = utils.NewBitSet(AuthMax)
		item.auth.UnmarshalBinary(byteAuth)
		err = item._initSubItemIds()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		_initMagnet(&item)
		items = append(items, &item)
	}
	if log.EnabledDebug() {
		ids := []ptype.CategoryID{}
		for _, item := range items {
			ids = append(ids, item.base.Id)
		}
		log.Debug("[category] loaded items: ", ids)
	}
	return items, nil
}

func addItem(params *NewCategoryParams) (*CategoryItem, error) {
	var newId ptype.CategoryID

	if params.CompareName {
		var c int
		err := db.QueryRow("select count(*) from pnas.category_items where parent_id=? and name=?", params.ParentId, params.Name).Scan(&c)
		if err != nil {
			return nil, err
		}
		if c > 0 {
			return nil, errors.New(fmt.Sprintf("existed name: %s", params.Name))
		}
	}

	if params.Auth.BitSet == nil {
		params.Auth = utils.NewBitSet(AuthMax)
	}
	byteAuth, err := params.Auth.MarshalBinary()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = db.QueryRow("call new_category(?, ?, ?, ?, ?, ?, ?, ?, @new_item_id)",
		params.TypeId,
		params.Name,
		params.Creator,
		byteAuth,
		params.ResourcePath,
		params.PosterPath,
		params.Introduce,
		params.ParentId).Scan(&newId)
	if err != nil {
		log.Warn(err)
		return nil, errors.WithStack(err)
	}
	if newId == NotExisted {
		return nil, errors.New("parent not existed")
	}
	return _loadItem(newId)
}

func (c *CategoryItem) _initSubItemIds() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	sql := `select id from pnas.category_items where parent_id=? order by type_id, name`
	rows, err := db.Query(sql, c.base.Id)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rows.Close()
	var subIds []ptype.CategoryID
	for rows.Next() {
		var itemId ptype.CategoryID
		err := rows.Scan(&itemId)
		if err != nil {
			log.Warn(err)
			continue
		}
		subIds = append(subIds, itemId)
	}
	c.subItemIds = subIds
	return nil
}

func (c *CategoryItem) deletedSubItem(subItemId ptype.CategoryID) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	index := slices.IndexFunc(c.subItemIds, func(sid ptype.CategoryID) bool { return subItemId == sid })
	if index == -1 {
		log.Warnf("[category] not found sub id: %d", subItemId)
		return
	}
	c.subItemIds = append(c.subItemIds[:index], c.subItemIds[index+1:]...)
}

func (c *CategoryItem) GetItemBaseInfo() BaseItem {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base
}

func (c *CategoryItem) GetId() ptype.CategoryID {
	return c.base.Id
}

func (c *CategoryItem) GetParentId() ptype.CategoryID {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.ParentId
}

func (c *CategoryItem) GetName() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.Name
}

func (c *CategoryItem) Rename(newName string) error {
	sql := "update pnas.category_items set name=? where id=?"
	_, err := db.Exec(sql, newName, c.base.Id)
	if err != nil {
		return errors.WithStack(err)
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.base.Name = newName
	c.base.UpdatedAt = time.Now()
	return nil
}

func (c *CategoryItem) GetOwner() ptype.UserID {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.Creator
}

func (c *CategoryItem) HasAndAuths(auths ...uint) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, i := range auths {
		if !c.auth.Test(i) {
			return false
		}
	}
	return true
}

func (c *CategoryItem) HasReadAuth(who ptype.UserID) bool {
	if who == ptype.AdminId {
		return true
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.base.Creator == who || c.auth.Test(AuthOtherRead) {
		return true
	}
	return false
}

func (c *CategoryItem) HasWriteAuth(who ptype.UserID) bool {
	if who == ptype.AdminId {
		return true
	}
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.base.Creator == who || c.auth.Test(AuthOtherWrite) {
		return true
	}
	return false
}

func (c *CategoryItem) GetType() prpc.CategoryItem_Type {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.TypeId
}

func (c *CategoryItem) IsDirectory() bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.TypeId == prpc.CategoryItem_Directory || c.base.TypeId == prpc.CategoryItem_Home
}

func (c *CategoryItem) GetSubItemIds() []ptype.CategoryID {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	ret := make([]ptype.CategoryID, len(c.subItemIds))
	copy(ret, c.subItemIds)
	return ret
}

func (c *CategoryItem) GetSubItemIdsInTypes(types *categoryTypes) []ptype.CategoryID {
	if types == nil || types.All() {
		return c.GetSubItemIds()
	}

	ids := types.GetTypeIds()
	idsStrBuilder := strings.Builder{}
	for _, id := range ids {
		if idsStrBuilder.Len() > 0 {
			idsStrBuilder.WriteString(",")
		}
		idsStrBuilder.WriteString(strconv.Itoa(int(id)))
	}
	sql := fmt.Sprintf(`select id from category_items where parent_id=? and type_id in (%s) order by type_id, name`, idsStrBuilder.String())
	rows, err := db.Query(sql, c.base.Id)
	if err != nil {
		log.Warn(err)
		return []ptype.CategoryID{}
	}
	defer rows.Close()
	var subIds []ptype.CategoryID
	for rows.Next() {
		var itemId ptype.CategoryID
		err := rows.Scan(&itemId)
		if err != nil {
			log.Warn(err)
			continue
		}
		subIds = append(subIds, itemId)
	}
	return subIds
}

func (c *CategoryItem) UpdatePosterPath(path string) error {
	sql := "update pnas.category_items set poster_path=? where id=?"
	_, err := db.Exec(sql, path, c.base.Id)
	if err == nil {
		c.mtx.Lock()
		defer c.mtx.Unlock()
		c.base.PosterPath = path
	} else {
		log.Warnf("itemId: %d, %v", c.base.Id, err)
	}
	return err
}

func (c *CategoryItem) GetPosterPath() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.PosterPath
}

func (c *CategoryItem) GetVideoId() ptype.VideoID {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.base.TypeId != prpc.CategoryItem_Video {
		return -1
	}
	vid, err := strconv.ParseInt(c.base.ResourcePath, 10, 64)
	if err != nil {
		return -1
	}
	return ptype.VideoID(vid)
}

func (c *CategoryItem) GetOther() OtherInfo {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.base.Other
}

func (c *CategoryItem) UpdateMagnetUri(magnetUri string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.base.Other.MagnetUri = magnetUri
}
