package user

import (
	"pnas/category"
	"pnas/log"
	"pnas/prpc"

	"github.com/pkg/errors"
)

const (
	rootDirectoryName = "magnet-shares"
)

type IMagnetSharesService interface {
	GetRootId() category.ID
	AddMagnetCategory(name string, parentId category.ID) error
	AddMagnetUri(params *AddMagnetUriParams) error
	QueryMagnetCategorys(params *QueryCategoryParams) ([]*category.CategoryItem, error)
	DelMagnetCategory(category.ID) error
}

type MagnetSharesService struct {
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

func (m *MagnetSharesService) AddMagnetCategory(name string, parentId category.ID) error {
	if !m.categoryService.IsRelationOf(m.GetRootId(), parentId) {
		return errors.New("isn't share directory")
	}
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
	Uri        string
	CategoryId category.ID
	Name       string
}

func (m *MagnetSharesService) AddMagnetUri(params *AddMagnetUriParams) error {
	_, err := m.categoryService.AddItem(&category.NewCategoryParams{
		ParentId:     params.CategoryId,
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
	if !m.categoryService.IsRelationOf(m.GetRootId(), params.ParentId) {
		return nil, errors.New("isn't share directory")
	}
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
	return m.queryMagnetCategorys(params)
}

func (m *MagnetSharesService) DelMagnetCategory(id category.ID) error {
	if !m.categoryService.IsRelationOf(m.GetRootId(), id) {
		return errors.New("isn't share directory")
	}
	return m.categoryService.DelItem(category.AdminId, id)
}
