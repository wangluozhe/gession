# Gession会话框架

[![Gitee link address](https://img.shields.io/badge/gitee-reference-red?logo=gitee&logoColor=red&labelColor=white)](https://gitee.com/leegene/gession)
[![Github link address](https://img.shields.io/badge/github-reference-blue?logo=github&logoColor=black&labelColor=white&color=black)](https://github.com/wangluozhe/gession)
[![Go Version](https://img.shields.io/badge/Go%20Version-1.15.6-blue?logo=go&logoColor=white&labelColor=gray)]()
[![Release Version](https://img.shields.io/badge/release-v1.0.0-blue)]()

1. Gession是一个会话框架。
2. 可用于HTTP/HTTPS网络连接的会话保持。
3. 内置对接Redis接口，连接保存Session到Redis。
4. 操作极其方便，易上手，源码易懂有注释。

# 安装
```shell script
$ go get -u github.com/wangluozhe/gession
```

# 依赖
[![Redigo v1.8.5](https://img.shields.io/badge/redigo-v1.8.5-blue?logo=go&logoColor=white)](https://github.com/gomodule/redigo)
[![go.uuid v1.2.0](https://img.shields.io/badge/go.uuid-v1.2.0-blue?logo=go&logoColor=white)](https://github.com/satori/go.uuid)

# 导入
```go
import(
    "github.com/wangluozhe/gession/session"
)
```

# 初始化Session管理器
Session管理器用于管理Session的各种操作，如：New/Get/Del。
Session管理器变量名：`session.Ssmgr`。

`Redis方式：`

```go
host := "host"				    // Redis地址
port := 6379				    // Redis端口
password := "password"		    // Redis密码
database := 1				    // Redis库
pool := session.NewRedisPool(host, int(port), password, int(database))		// 创建Redis连接池
expire := 1000				    // Session过期时间
session.Init(pool, expire)	    // 初始化Session管理器
Ssmgr := session.Ssmgr		    // 全局可用，可以直接操作session.Ssmgr也可以赋值后在操作，Ssmgr是Session组管理器，用来创建/读取/删除Session
```

`内存方式：`
```go
session.Init()                  // 初始化Session管理器
Ssmgr := session.Ssmgr		    // 全局可用，可以直接操作session.Ssmgr也可以赋值后在操作，Ssmgr是Session组管理器，用来创建/读取/删除Session
```

### 创建Session

```go
ss := Ssmgr.New()			    // 创建Session
```

### 读取Session
```go
ss,err := Ssmgr.Get(sessionId)  // 读取Session，如果内存中没有此Session会从Redis中读取，都没有返回nil
if err != nil{
    fmt.Println(err)
}
```

### 删除Session
```go
ss,err := Ssmgr.Del(sessionId)  // 删除Session，内存和Redis中的此Session都会被删除
if err != nil{
    fmt.Println(err)
}
```

# Session操作
Session操作只有四种操作，简单方便，存储的是map类型，有：Set/Get/Del/Save。

### 设置Session值
```go
ss.Set(key,value)               // 设置Session中的值
// 或
isSuccess := ss.Set(key,value)  // 返回一个是否设置成功的bool值
```

### 获取Session值
```go
result := ss.Get(key)           // 获取Session中的值
```

### 删除Session值
```go
ss.Del(key)                     // 删除Session中的值
// 或
isSuccess := ss.Del(key)        // 返回一个是否删除成功的bool值
```

### 保存Session到Redis中
```go
isSuccess,err := ss.Save()      // 保存此Session到Redis中
if err != nil{
    fmt.Println(err)
}
```