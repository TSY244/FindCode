package msg

import (
	"ScanIDOR/pkg/set"
	"sync"
)

type Msg struct {
	sync.RWMutex
	data *set.Set[string]
}

func NewMsg() *Msg {
	return &Msg{
		RWMutex: sync.RWMutex{},
		data:    set.NewSet[string](),
	}
}

func (m *Msg) AddFinishedTask(taskId string) {
	m.Lock()
	defer m.Unlock()
	m.data.Add(taskId)
}

func (m *Msg) In(taskId string) bool {
	m.RLock()
	defer m.RUnlock()
	return m.data.Contains(taskId)
}
