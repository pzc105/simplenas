package category

import (
	"pnas/ptype"
)

type IService interface {
	Init()
	AddItem(params *NewCategoryParams) (*CategoryItem, error)
	GetItem(querier ptype.UserID, itemId ptype.CategoryID) (*CategoryItem, error)
	GetItemsByParent(*GetItemsByParentParams) ([]*CategoryItem, error)
	GetItems(querier ptype.UserID, itemIds ...ptype.CategoryID) ([]*CategoryItem, error)
	GetItemByName(querier ptype.UserID, parentId ptype.CategoryID, name string) (*CategoryItem, error)
	DelItem(deleter ptype.UserID, itemId ptype.CategoryID) error
	IsRelationOf(itemId ptype.CategoryID, parentId ptype.CategoryID) bool
	Search(params *SearchParams) ([]*CategoryItem, error)
	SearchRows(params *SearchParams) (int, error)
}

type GetItemsByParentParams struct {
	Querier  ptype.UserID
	ParentId ptype.CategoryID
	PageNum  int32
	Rows     int32
}
