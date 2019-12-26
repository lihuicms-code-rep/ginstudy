package session

import (
	"errors"
	uuid2 "github.com/satori/go.uuid"
	"sync"
)

//内存session管理者
type MemorySessionMgr struct {
	sessions map[string] *MemorySession    //管理的所有Session
	rwLock sync.RWMutex                    //锁
}

func NewMemorySessinMgr() *MemorySessionMgr {
	return &MemorySessionMgr{
		sessions:make(map[string]*MemorySession),
	}
}

func (msMgr *MemorySessionMgr) Init(addr string, options ...string) error {
     return nil
}

func (msMgr *MemorySessionMgr) CreateSession() (ISession, error) {
     msMgr.rwLock.Lock()
     defer msMgr.rwLock.Unlock()

     //以uuid作为sessionID
     uuid := uuid2.NewV4()
     sessionID := uuid.String()
     session := NewMemorySession(sessionID)
     msMgr.sessions[sessionID] = session
     return session, nil
}


func (msMgr *MemorySessionMgr) Get(sessionID string) (ISession, error) {
     msMgr.rwLock.RLock()
     defer msMgr.rwLock.RUnlock()

     session, ok := msMgr.sessions[sessionID]
     if ok {
     	return session, nil
	 }

     return nil, errors.New("session not exist")
}