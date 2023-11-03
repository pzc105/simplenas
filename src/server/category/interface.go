package category

type IService interface {
	Init()
	AddItem(params *NewCategoryParams) (*CategoryItem, error)
	GetItem(querier int64, itemId ID) (*CategoryItem, error)
	GetItems(querier int64, itemIds ...ID) ([]*CategoryItem, error)
	GetItemByName(querier int64, parentId ID, name string) (*CategoryItem, error)
	DelItem(deleter int64, itemId ID) error
	IsRelationOf(parentId ID, itemId ID) bool
	Search(params *SearchParams) ([]*CategoryItem, error)
}
