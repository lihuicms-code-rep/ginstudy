package session

import (
	"errors"
	"sync"
)

//内存实现session
type MemorySession struct {
	sessionID string             //sessionID
	data map[string]interface{}  //具体存储信息
	rwLock sync.RWMutex          //锁
}

//构造函数
func NewMemorySession(sID string) *MemorySession {
	return &MemorySession{
		sessionID:sID,
		data:make(map[string]interface{}),
	}
}


func (ms *MemorySession) Set (key string, value interface{}) error {
	ms.rwLock.Lock()
	defer ms.rwLock.Unlock()

	ms.data[key] = value
	return nil
}

func (ms *MemorySession) Get(key string) (interface{}, error) {
	ms.rwLock.Lock()
	defer ms.rwLock.Unlock()

	value, ok := ms.data[key]
	if ok {
		return value, nil
	}

	return nil, errors.New("key not exist in session")
}

func (ms *MemorySession) Del(key string) error {
	ms.rwLock.Lock()
	defer ms.rwLock.Unlock()
	delete(ms.data, key)
    return nil
}


func (ms *MemorySession) Save() error {
    return nil
}