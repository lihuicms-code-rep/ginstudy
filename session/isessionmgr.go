package session

//session管理抽象层
type ISessionMgr interface {
	Init(addr string, options ...string) error
	CreateSession() (ISession, error)
	Get(sessionID string) (ISession, error)
}
