package session

var Ssmgr *RedisSessionMgr

func Init(options ...interface{}) {
	if len(options) == 0 {
		Ssmgr = InitSessionMgr()
	} else {
		Ssmgr = InitSessionMgr(options...)
	}
}