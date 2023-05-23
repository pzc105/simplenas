package category

import (
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/setting"
	"pnas/utils"
	"testing"
)

func init() {
	setting.Init("../")
	log.Init()
	db.Init()
}

func TestNewItem(t *testing.T) {
	var m Manager
	m.Init()
	auth := utils.NewBitSet(AuthMax)
	auth.Set(AuthOtherWrite)
	otherId := int64(124)
	params := &NewCategoryParams{
		ParentId:     2,
		Creator:      123,
		TypeId:       prpc.CategoryItem_Directory,
		Name:         "123",
		ResourcePath: "resource",
		Introduce:    "introduce",
		Auth:         auth,
	}
	item1, err := m.NewItem(params)
	if err != nil {
		t.Errorf("failed to create: %v", err)
		return
	}
	item2, err := _loadItem(item1.base.Id)
	if err != nil {
		t.Errorf("failed to load: %v", err)
		return
	}

	if item1.base.Id != item2.base.Id {
		t.Error("id not equal")
	}
	if item1.base.ParentId != item2.base.ParentId {
		t.Error("parent id not equal")
	}
	if item1.base.ResourcePath != item2.base.ResourcePath || item1.base.ResourcePath != params.ResourcePath {
		t.Error("resource path not equal")
	}
	if item1.base.Introduce != item2.base.Introduce || item1.base.Introduce != params.Introduce {
		t.Error("introduce not equal")
	}
	if item1.HasWriteAuth(otherId) != item2.HasWriteAuth(otherId) || !item1.HasWriteAuth(otherId) {
		t.Errorf("write auth not equal %s %s", item1.auth.String(), item2.auth.String())
	}
	if item1.HasReadAuth(otherId) != item2.HasReadAuth(otherId) || item1.HasReadAuth(otherId) {
		t.Error("read auth not equal")
	}
	err = m.DelItem(item1.base.Id)
	if err != nil {
		t.Errorf("failed to del: %v", err)
		return
	}
	if eitem, _ := m.GetItem(item1.base.Id); eitem != nil {
		t.Errorf("not be deleted item cache: %v", err)
	}
	if eitem, _ := _loadItem(item1.base.Id); eitem != nil {
		t.Errorf("not be deleted item: %v", err)
	}
}
