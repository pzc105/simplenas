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
	GetMagnetRootId() category.ID
	AddMagnetCategory(params *AddMagnetCategoryParams) error
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

func (m *MagnetSharesService) GetMagnetRootId() category.ID {
	return m.rootId
}

type AddMagnetCategoryParams struct {
	ParentId  category.ID
	Name      string
	Introduce string
}

func (m *MagnetSharesService) AddMagnetCategory(params *AddMagnetCategoryParams) error {
	if !m.categoryService.IsRelationOf(m.GetMagnetRootId(), params.ParentId) {
		return errors.New("isn't share directory")
	}
	_, err := m.categoryService.AddItem(&category.NewCategoryParams{
		ParentId:    params.ParentId,
		Creator:     category.AdminId,
		TypeId:      prpc.CategoryItem_Directory,
		Name:        params.Name,
		Introduce:   params.Introduce,
		CompareName: true,
	})
	return err
}

type AddMagnetUriParams struct {
	Uri        string
	CategoryId category.ID
	Name       string
	Introduce  string
}

func (m *MagnetSharesService) AddMagnetUri(params *AddMagnetUriParams) error {
	_, err := m.categoryService.AddItem(&category.NewCategoryParams{
		ParentId:  params.CategoryId,
		Creator:   category.AdminId,
		TypeId:    prpc.CategoryItem_Other,
		Name:      params.Name,
		Other:     params.Uri,
		Introduce: params.Introduce,
	})
	return err
}

type QueryCategoryParams struct {
	ParentId     category.ID
	CategoryName string
}

func (m *MagnetSharesService) queryMagnetCategorys(params *QueryCategoryParams) ([]*category.CategoryItem, error) {
	if !m.categoryService.IsRelationOf(m.GetMagnetRootId(), params.ParentId) {
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
	if !m.categoryService.IsRelationOf(m.GetMagnetRootId(), id) {
		return errors.New("isn't share directory")
	}
	return m.categoryService.DelItem(category.AdminId, id)
}
