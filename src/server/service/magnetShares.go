package service

import (
	"pnas/category"
	"pnas/log"
	"pnas/prpc"
)

const (
	rootDirectoryName = "magnet-shares"
)

type MagnetSharesService struct {
	rootId          category.ID
	categoryService category.Service
}

func (m *MagnetSharesService) Init(ser category.Service) {
	m.categoryService = ser
	item, err := ser.GetItem(category.AdminId, category.RootId)
	if err != nil {
		log.Errorf("failed to load root: %v", err)
		return
	}
	found := false
	for _, sid := range item.GetSubItemIds() {
		item, err := ser.GetItem(category.AdminId, sid)
		if err != nil {
			continue
		}
		ii := item.GetItemInfo()
		if ii.Name == rootDirectoryName {
			m.rootId = ii.Id
			found = true
			break
		}
	}
	if !found {
		_, err := ser.AddItem(&category.NewCategoryParams{
			ParentId: category.RootId,
			Creator:  category.AdminId,
			TypeId:   prpc.CategoryItem_Directory,
			Name:     rootDirectoryName,
		})
		if err != nil {
			log.Errorf("failed to create magnet rootDirectory: %v", err)
		}
	}
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
	Name         string
	ResourcePath string
}

func (m *MagnetSharesService) AddMagnetUri(params *AddMagnetUriParams) error {
	_, err := m.categoryService.AddItem(&category.NewCategoryParams{
		ParentId:     params.ParentId,
		Creator:      category.AdminId,
		TypeId:       prpc.CategoryItem_Other,
		Name:         params.Name,
		ResourcePath: params.ResourcePath,
	})
	return err
}

type QueryCategoryParams struct {
	ParentId category.ID
}

func (m *MagnetSharesService) queryCategory(params *QueryCategoryParams) ([]*category.CategoryItem, error) {
	pitem, err := m.categoryService.GetItem(category.AdminId, params.ParentId)
	if err != nil {
		return nil, err
	}
	ret, err := m.categoryService.GetItems(category.AdminId, pitem.GetSubItemIds()...)
	return ret, err
}

func (m *MagnetSharesService) QueryCategory(params *QueryCategoryParams) ([]*category.CategoryItem, error) {
	if params.ParentId <= 0 {
		params.ParentId = category.AdminId
	}
	return m.queryCategory(params)
}
