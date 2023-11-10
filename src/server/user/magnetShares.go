package user

import (
	"pnas/category"
	"pnas/log"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/utils"

	"github.com/pkg/errors"
)

const (
	rootDirectoryName = "magnet-shares"
)

type IMagnetSharesService interface {
	GetMagnetRootId() ptype.CategoryID
	AddMagnetCategory(params *AddMagnetCategoryParams) (ptype.CategoryID, error)
	AddMagnetUri(params *AddMagnetUriParams) error
	QueryMagnetCategorys(params *QueryCategoryParams) ([]*category.CategoryItem, error)
	DelMagnetCategory(ptype.UserID, ptype.CategoryID) error
}

type MagnetSharesService struct {
	rootId          ptype.CategoryID
	categoryService category.IService
}

func (m *MagnetSharesService) Init(ser category.IService) {
	m.categoryService = ser
	item, err := ser.GetItemByName(ptype.AdminId, category.RootId, rootDirectoryName)
	if err != nil {
		item, err = ser.AddItem(&category.NewCategoryParams{
			ParentId: category.RootId,
			Creator:  ptype.AdminId,
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

func (m *MagnetSharesService) GetMagnetRootId() ptype.CategoryID {
	return m.rootId
}

type AddMagnetCategoryParams struct {
	ParentId  ptype.CategoryID
	Name      string
	Introduce string
	Creator   ptype.UserID
}

func (m *MagnetSharesService) AddMagnetCategory(params *AddMagnetCategoryParams) (ptype.CategoryID, error) {
	if !m.categoryService.IsRelationOf(params.ParentId, m.GetMagnetRootId()) {
		return -1, errors.New("isn't share directory")
	}
	sudo := false
	if params.ParentId == m.rootId {
		sudo = true
	}
	item, err := m.categoryService.AddItem(&category.NewCategoryParams{
		ParentId:    params.ParentId,
		Creator:     params.Creator,
		TypeId:      prpc.CategoryItem_Directory,
		Name:        params.Name,
		Introduce:   params.Introduce,
		Auth:        utils.NewBitSet(category.AuthMax, category.AuthOtherRead),
		CompareName: true,
		Sudo:        sudo,
	})
	if err != nil {
		return -1, nil
	}
	return item.GetItemInfo().Id, err
}

type AddMagnetUriParams struct {
	Uri        string
	CategoryId ptype.CategoryID
	Name       string
	Introduce  string
	Creator    ptype.UserID
}

func (m *MagnetSharesService) AddMagnetUri(params *AddMagnetUriParams) error {
	_, err := m.categoryService.AddItem(&category.NewCategoryParams{
		ParentId:  params.CategoryId,
		Creator:   params.Creator,
		TypeId:    prpc.CategoryItem_Other,
		Name:      params.Name,
		Other:     params.Uri,
		Introduce: params.Introduce,
		Auth:      utils.NewBitSet(category.AuthMax, category.AuthOtherRead),
	})
	return err
}

type QueryCategoryParams struct {
	ParentId     ptype.CategoryID
	CategoryName string
	PageNum      int32
	Rows         int32
}

func (m *MagnetSharesService) queryMagnetCategorys(params *QueryCategoryParams) ([]*category.CategoryItem, error) {
	if !m.categoryService.IsRelationOf(params.ParentId, m.GetMagnetRootId()) {
		return nil, errors.New("isn't share directory")
	}
	if len(params.CategoryName) != 0 {
		item, err := m.categoryService.GetItemByName(ptype.AdminId, params.ParentId, params.CategoryName)
		if err != nil {
			return nil, err
		}
		return []*category.CategoryItem{item}, nil
	}
	ret, err := m.categoryService.GetItemsByParent(&category.GetItemsByParentParams{
		Querier:  ptype.AdminId,
		ParentId: params.ParentId,
		PageNum:  params.PageNum,
		Rows:     params.Rows,
	})
	return ret, err
}

func (m *MagnetSharesService) QueryMagnetCategorys(params *QueryCategoryParams) ([]*category.CategoryItem, error) {
	return m.queryMagnetCategorys(params)
}

func (m *MagnetSharesService) DelMagnetCategory(deletor ptype.UserID, id ptype.CategoryID) error {
	if !m.categoryService.IsRelationOf(id, m.GetMagnetRootId()) {
		return errors.New("isn't share directory")
	}
	return m.categoryService.DelItem(deletor, id)
}
