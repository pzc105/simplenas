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
		_, err := ser.NewItem(&category.NewCategoryParams{
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

func (m *MagnetSharesService) NewRootCategory(name string) error {
	_, err := m.categoryService.NewItem(&category.NewCategoryParams{
		ParentId: m.rootId,
		Creator:  category.AdminId,
		TypeId:   prpc.CategoryItem_Directory,
		Name:     name,
	})
	return err
}

func (m *MagnetSharesService) NewCategory(name string, parentId category.ID) error {
	_, err := m.categoryService.NewItem(&category.NewCategoryParams{
		ParentId: parentId,
		Creator:  category.AdminId,
		TypeId:   prpc.CategoryItem_Directory,
		Name:     name,
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
	_, err := m.categoryService.NewItem(&category.NewCategoryParams{
		ParentId:     params.ParentId,
		Creator:      category.AdminId,
		TypeId:       prpc.CategoryItem_Other,
		Name:         params.Name,
		ResourcePath: params.ResourcePath,
	})
	return err
}
