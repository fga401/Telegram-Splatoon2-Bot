package adapter

type Manager interface {
	Index(adapter Adapter) []int
	Add(adapter Adapter)
}

type managerImpl struct {
	argMap map[string][]int
}

func (m *managerImpl) Add(adapter Adapter) {
	nArgs:= adapter.ArgNum()
	for _, idx := range m.argMap {
		for i:=0;i< len(idx);i++{
			idx[i] += nArgs
		}
	}
	argIdx := make([]int, nArgs)
	for i := 0; i < len(argIdx); i++ {
		argIdx[i] = i
	}
	m.argMap[adapter.ID()] = argIdx
}

func (m *managerImpl) Index(adapter Adapter) []int {
	return m.argMap[adapter.ID()]
}

func NewManager() Manager {
	return &managerImpl{
		argMap: make(map[string][]int),
	}
}
