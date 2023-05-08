package models

import (
	"sync"
)

func NewSyncMap() *SyncMap {
	return &SyncMap{M: make(map[int64]User)}
}

type SyncMap struct {
	sync.RWMutex
	M map[int64]User
}

func (sm *SyncMap) SafeUpdateUser(user User) {
	sm.Lock()
	sm.M[user.TGID] = user
	sm.Unlock()
}
