package session

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"sync"
)

// Session管理器结构体
type RedisSessionMgr struct {
	pool        *redis.Pool
	expire      int
	rwlock      sync.RWMutex
	sessionMaps map[string]*RedisSession
}

// 初始化Session管理器
func InitSessionMgr(options ...interface{}) *RedisSessionMgr {
	if len(options) == 0 {
		return &RedisSessionMgr{
			sessionMaps: make(map[string]*RedisSession, 0),
		}
	} else if len(options) == 1 {
		return &RedisSessionMgr{
			pool:        options[0].(*redis.Pool),
			sessionMaps: make(map[string]*RedisSession, 0),
		}
	}
	return &RedisSessionMgr{
		pool:        options[0].(*redis.Pool),
		expire:      options[1].(int),
		sessionMaps: make(map[string]*RedisSession, 0),
	}
}

// 新建Session
func (this *RedisSessionMgr) New(options ...string) *RedisSession {
	var sessionId string
	if len(options) == 0 {
		sessionId = uuid.NewV4().String()
	} else {
		sessionId = options[0]
	}
	if this.pool == nil {
		this.sessionMaps[sessionId] = initSession(sessionId)
	} else if this.expire == 0 {
		this.sessionMaps[sessionId] = initSession(sessionId, this.pool)
	} else {
		this.sessionMaps[sessionId] = initSession(sessionId, this.pool, this.expire)
	}
	return this.sessionMaps[sessionId]
}

// 读取Session
func (this *RedisSessionMgr) Get(sessionId string) (*RedisSession,error) {
	this.rwlock.Lock()
	defer this.rwlock.Unlock()
	session := this.sessionMaps[sessionId]
	if this.pool != nil && session == nil {
		session,err := this.loadFromRedis(sessionId)
		if err != nil{
			return nil,err
		}
		this.sessionMaps[sessionId] = session
		return session,nil
	}
	return session,nil
}

// 从Redis中加载Session
func (this *RedisSessionMgr) loadFromRedis(sessionId string) (*RedisSession,error) {
	conn := this.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", sessionId)
	if err != nil{
		return nil,err
	}
	if reply == nil {
		return nil,errors.New("Redis没有这个SessionId的数据")
	}
	result, err := redis.String(reply, err)
	session := this.New(sessionId)
	session.Flag = SessionFlagModify
	err = json.Unmarshal([]byte(result), &(session.SessionMap))
	if err != nil {
		return nil,err
	}
	return session,nil
}

// 从Redis中删除Session
func (this *RedisSessionMgr) Del(sessionId string) (bool,error) {
	delete(this.sessionMaps, sessionId)
	if this.pool != nil {
		conn := this.pool.Get()
		defer conn.Close()
		_, err := conn.Do("DEL", sessionId)
		if err != nil {
			return false,err
		}
	}
	return true,nil
}
