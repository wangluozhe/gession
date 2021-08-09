package session

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"sync"
)

// RedisSession结构体
type RedisSession struct {
	SessionId  string
	pool       *redis.Pool
	Expire     int
	SessionMap map[string]interface{}
	rwlock     sync.RWMutex
	Flag       int
}

const (
	SessionFlagNone   = iota // 无数据状态
	SessionFlagModify        // 有数据状态
)

// 初始化一个Session
func initSession(sessionId string, options ...interface{}) *RedisSession {
	if len(options) == 0 {
		return &RedisSession{
			SessionId:  sessionId,
			SessionMap: make(map[string]interface{}, 0),
			Flag:       SessionFlagNone,
		}
	} else if len(options) == 1 {
		return &RedisSession{
			SessionId:  sessionId,
			pool:       options[0].(*redis.Pool),
			SessionMap: make(map[string]interface{}, 0),
			Flag:       SessionFlagNone,
		}
	}
	return &RedisSession{
		SessionId:  sessionId,
		pool:       options[0].(*redis.Pool),
		Expire:     options[1].(int),
		SessionMap: make(map[string]interface{}, 0),
		Flag:       SessionFlagNone,
	}
}

// 设置Session内容
func (this *RedisSession) Set(key string, value interface{}) bool {
	this.rwlock.Lock()
	defer this.rwlock.Unlock()
	if this.Flag != SessionFlagModify {
		this.Flag = SessionFlagModify
	}
	this.SessionMap[key] = value
	return true
}

// 获取Session内容
func (this *RedisSession) Get(key string) interface{} {
	this.rwlock.Lock()
	defer this.rwlock.Unlock()
	if this.Flag != SessionFlagModify {
		return nil
	}
	result, ok := this.SessionMap[key]
	if !ok {
		return nil
	}
	return result
}

// 删除Session内容
func (this *RedisSession) Del(key string) bool {
	this.rwlock.Lock()
	defer this.rwlock.Unlock()
	if this.Flag != SessionFlagModify {
		return false
	}
	if len(this.SessionMap) == 0 {
		this.Flag = SessionFlagNone
	}
	delete(this.SessionMap, key)
	return true
}

// 保存Session内容
func (this *RedisSession) Save(expire ...int) (bool,error) {
	this.rwlock.Lock()
	defer this.rwlock.Unlock()
	if this.Flag != SessionFlagModify {
		return false,errors.New("当前Session没有数据")
	}
	if this.pool == nil {
		return false,errors.New("当前Session未连接Redis")
	}
	conn := this.pool.Get()
	defer conn.Close()
	data, err := json.Marshal(this.SessionMap)
	if err != nil {
		return false,err
	}
	_, err = conn.Do("SET", this.SessionId, string(data))
	if err != nil{
		return false,err
	}
	if len(expire) > 0 {
		_, err = conn.Do("EXPIRE", this.SessionId, expire[0])
		if err != nil{
			return false,err
		}
		return true,err
	}
	_, err = conn.Do("EXPIRE", this.SessionId, this.Expire)
	if err != nil{
		return false,err
	}
	return true,err
}
