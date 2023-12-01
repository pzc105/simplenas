package category

import (
	"pnas/prpc"
	"pnas/ptype"
	"pnas/utils"
)

type IService interface {
	Init()
	AddItem(params *NewCategoryParams) (*CategoryItem, error)
	GetItem(querier ptype.UserID, itemId ptype.CategoryID) (*CategoryItem, error)
	RefreshItem(itemId ptype.CategoryID) error
	GetItemsByParent(*GetItemsByParentParams) ([]*CategoryItem, error)
	GetItems(querier ptype.UserID, itemIds ...ptype.CategoryID) ([]*CategoryItem, error)
	GetItemByName(querier ptype.UserID, parentId ptype.CategoryID, name string) (*CategoryItem, error)
	DelItem(deleter ptype.UserID, itemId ptype.CategoryID) error
	IsRelationOf(itemId ptype.CategoryID, parentId ptype.CategoryID) bool
	Search(params *SearchParams) ([]*CategoryItem, error)
	SearchRows(params *SearchParams) (int, error)
}

type NewCategoryParams struct {
	ParentId     ptype.CategoryID
	Creator      ptype.UserID
	TypeId       prpc.CategoryItem_Type
	Name         string
	ResourcePath string
	PosterPath   string
	Introduce    string
	Other        OtherInfo
	Auth         utils.AuthBitSet
	CompareName  bool
	Sudo         bool
}

type GetItemsByParentParams struct {
	Querier  ptype.UserID
	ParentId ptype.CategoryID
	PageNum  int32
	Rows     int32
}

type SearchParams struct {
	Querier         ptype.UserID
	RootId          ptype.CategoryID
	ExistedWords    []string
	NotExistedWords []string
	PageNum         int32
	Rows            int32
}
