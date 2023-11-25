package user

import (
	"fmt"
	"pnas/bt"
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
	AddMagnetUriByTorrent(params *AddMagnetUriParams) error
	QueryMagnetCategorys(params *QueryCategoryParams) ([]*category.CategoryItem, error)
	DelMagnetCategory(ptype.UserID, ptype.CategoryID) error
}

type MagnetSharesService struct {
	rootId          ptype.CategoryID
	categoryService category.IService
	userTorrents    bt.UserTorrents
}

func (m *MagnetSharesService) Init(cser category.IService, ubts bt.UserTorrents) {
	m.categoryService = cser
	m.userTorrents = ubts
	item, err := cser.GetItemByName(ptype.AdminId, category.RootId, rootDirectoryName)
	if err != nil {
		item, err = cser.AddItem(&category.NewCategoryParams{
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

	m.rootId = item.GetItemBaseInfo().Id
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
	return item.GetItemBaseInfo().Id, err
}

type AddMagnetUriParams struct {
	CategoryId ptype.CategoryID
	Name       string
	Introduce  string
	Creator    ptype.UserID
	Uri        string
	Torrent    []byte
}

func (m *MagnetSharesService) AddMagnetUri(params *AddMagnetUriParams) error {
	t, _ := m.userTorrents.NewTorrentByMagnet(params.Uri)
	if t == nil {
		return errors.New("failed to new torrent")
	}
	_, err := m.categoryService.AddItem(&category.NewCategoryParams{
		ParentId:     params.CategoryId,
		Creator:      params.Creator,
		TypeId:       prpc.CategoryItem_MagnetUri,
		ResourcePath: fmt.Sprint(t.GetId()),
		Other: category.OtherInfo{
			MagnetUri: params.Uri,
		},
		Name:        params.Name,
		Introduce:   params.Introduce,
		CompareName: true,
		Auth:        utils.NewBitSet(category.AuthMax, category.AuthOtherRead),
	})
	return err
}

func (m *MagnetSharesService) AddMagnetUriByTorrent(params *AddMagnetUriParams) error {
	rsp, err := m.userTorrents.GetMagnetUri(&bt.GetMagnetUriParams{
		Req: &prpc.GetMagnetUriReq{
			Type:    prpc.GetMagnetUriReq_Torrent,
			Content: params.Torrent,
		},
	})
	if err != nil {
		return err
	}
	params.Uri = rsp.MagnetUri
	return m.AddMagnetUri(params)
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
