package user

import (
	"pnas/category"
	"pnas/log"
	"pnas/prpc"
)

const (
	rootDirectoryName = "magnet-shares"
)

type IMagnetSharesService interface {
	GetRootId() category.ID
	AddMagnetRootCategory(name string) error
	AddMagnetCategory(name string, parentId category.ID) error
	AddMagnetUri(params *AddMagnetUriParams) error
	QueryMagnetCategorys(params *QueryCategoryParams) ([]*category.CategoryItem, error)
}

type MagnetSharesService struct {
	IMagnetSharesService
	rootId          category.ID
	categoryService category.Service
}

func (m *MagnetSharesService) Init(ser category.Service) {
	m.categoryService = ser
	item, err := ser.GetItemByName(category.AdminId, category.RootId, rootDirectoryName)
	if err != nil {
		item, err = ser.AddItem(&category.NewCategoryParams{
			ParentId: category.RootId,
			Creator:  category.AdminId,
			TypeId:   prpc.CategoryItem_Directory,
			Name:     rootDirectoryName,
		})
		if err != nil {
			log.Errorf("failed to create magnet rootDirectory: %v", err)
			return
		}
	}

	m.rootId = item.GetItemInfo().Id
}

func (m *MagnetSharesService) GetRootId() category.ID {
	return m.rootId
}

func (m *MagnetSharesService) AddRootCategory(name string) error {
	_, err := m.categoryService.AddItem(&category.NewCategoryParams{
		ParentId:    m.rootId,
		Creator:     category.AdminId,
		TypeId:      prpc.CategoryItem_Directory,
		Name:        name,
		CompareName: true,
	})
	return err
}

func (m *MagnetSharesService) AddCategory(name string, parentId category.ID) error {
	_, err := m.categoryService.AddItem(&category.NewCategoryParams{
		ParentId:    parentId,
		Creator:     category.AdminId,
		TypeId:      prpc.CategoryItem_Directory,
		Name:        name,
		CompareName: true,
	})
	return err
}

type AddMagnetUriParams struct {
	Uri          string
	ParentId     category.ID
	CategoryName string
	Name         string
}

func (m *MagnetSharesService) AddMagnetUri(params *AddMagnetUriParams) error {
	parentId := params.ParentId
	if parentId <= 0 {
		params.ParentId = category.AdminId
	}
	_, err := m.categoryService.GetItemByName(category.AdminId, parentId, params.CategoryName)
	if err != nil {
		err = m.AddCategory(params.CategoryName, parentId)
		if err != nil {
			return err
		}
	}
	_, err = m.categoryService.AddItem(&category.NewCategoryParams{
		ParentId:     params.ParentId,
		Creator:      category.AdminId,
		TypeId:       prpc.CategoryItem_Other,
		Name:         params.Name,
		ResourcePath: params.Uri,
	})
	return err
}

type QueryCategoryParams struct {
	ParentId     category.ID
	CategoryName string
}

func (m *MagnetSharesService) queryMagnetCategorys(params *QueryCategoryParams) ([]*category.CategoryItem, error) {
	if len(params.CategoryName) != 0 {
		item, err := m.categoryService.GetItemByName(category.AdminId, params.ParentId, params.CategoryName)
		if err != nil {
			return nil, err
		}
		return []*category.CategoryItem{item}, nil
	}
	pitem, err := m.categoryService.GetItem(category.AdminId, params.ParentId)
	if err != nil {
		return nil, err
	}
	ret, err := m.categoryService.GetItems(category.AdminId, pitem.GetSubItemIds()...)
	return ret, err
}

func (m *MagnetSharesService) QueryMagnetCategorys(params *QueryCategoryParams) ([]*category.CategoryItem, error) {
	if params.ParentId <= 0 {
		params.ParentId = category.AdminId
	}
	return m.queryMagnetCategorys(params)
}
