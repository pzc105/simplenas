package utils

import "testing"

func TestIdPool(t *testing.T) {
	var idPool IdPool
	idPool.Init()

	id1 := idPool.NewId()
	id2 := idPool.NewId()

	if id1+1 != id2 {
		t.Errorf("id1 + 1 != id2")
		return
	}

	idPool.ReleaseId(id1)
	id3 := idPool.NewId()
	if id3 != id1 {
		t.Errorf("id3 != id1")
		return
	}

	id4 := id2 + 1
	idPool.Allocated(id4)
	id5 := idPool.NewId()
	if id4+1 != id5 {
		t.Errorf("id4 != id5")
		return
	}
}
