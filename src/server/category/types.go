package category

import "pnas/prpc"

var (
	AllTypes categoryTypes
)

func init() {
	AllTypes.all = true
}

type categoryTypes struct {
	all   bool
	types map[prpc.CategoryItem_Type]bool
}

func NewCategoryTypes(types ...prpc.CategoryItem_Type) *categoryTypes {
	if len(types) == 0 {
		return &AllTypes
	}
	ret := &categoryTypes{
		all:   false,
		types: map[prpc.CategoryItem_Type]bool{},
	}
	for _, st := range types {
		ret.types[st] = true
	}
	return ret
}

func (t *categoryTypes) Add(st prpc.CategoryItem_Type) {
	t.types[st] = true
}

func (t *categoryTypes) Del(st prpc.CategoryItem_Type) {
	delete(t.types, st)
}

func (t *categoryTypes) Has(st prpc.CategoryItem_Type) bool {
	if f, ok := t.types[st]; ok && f {
		return true
	}
	return false
}

func (t *categoryTypes) All() bool {
	return t.all
}

func (t *categoryTypes) GetTypeIds() []prpc.CategoryItem_Type {
	ret := []prpc.CategoryItem_Type{}
	for k, v := range t.types {
		if v {
			ret = append(ret, k)
		}
	}
	return ret
}
