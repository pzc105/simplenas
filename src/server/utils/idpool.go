package utils

import "sync"

type IdPool struct {
	idleIds      []int64
	mtx          sync.Mutex
	idIndex      int64
	allocatedIds map[int64]bool
}

func (p *IdPool) Init() {
	p.idIndex = 0
	p.allocatedIds = make(map[int64]bool)
}

func (p *IdPool) NewId() int64 {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if len(p.idleIds) > 0 {
		ret := p.idleIds[0]
		p.idleIds = p.idleIds[1:]
		p.allocatedIds[ret] = true
		return ret
	}
	ret := p.idIndex
	for b, ok := p.allocatedIds[ret]; ok && b; b, ok = p.allocatedIds[ret] {
		p.idIndex += 1
		ret = p.idIndex
	}
	p.allocatedIds[ret] = true
	p.idIndex += 1
	return ret
}

func (p *IdPool) Allocated(id int64) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.allocatedIds[id] = true
}

func (p *IdPool) ReleaseId(id int64) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	_, ok := p.allocatedIds[id]
	if !ok {
		return
	}
	delete(p.allocatedIds, id)

	p.idleIds = append(p.idleIds, id)
}

func (p *IdPool) GetAllocatedIds() []int64 {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	keys := make([]int64, 0, len(p.allocatedIds))
	for k := range p.allocatedIds {
		keys = append(keys, k)
	}
	return keys
}
