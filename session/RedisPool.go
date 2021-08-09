package session

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

// 创建Redis连接池
func NewRedisPool(host string, port int, password string, database int) *redis.Pool {
	return &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			var (
				c   redis.Conn
				err error
			)
			address := host + ":" + strconv.Itoa(port)
			if password == "" {
				c, err = redis.Dial("tcp", address)
				if err != nil {
					panic(err)
					return nil, err
				}
			} else {
				dialpassword := redis.DialPassword(password)
				c, err = redis.Dial("tcp", address, dialpassword)
				if err != nil {
					panic(err)
					return nil, err
				}
			}
			c.Do("SELECT", database)
			return c, err
		},
		MaxIdle:         1000,              // 最大空闲连接数
		MaxActive:       2000,              // 最大连接数
		IdleTimeout:     180 * time.Second, // 空闲连接过期时效
		Wait:            true,              // 开启等待状态
		MaxConnLifetime: 0,
	}
}
