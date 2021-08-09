# Gession Session Framework

[![Gitee link address](https://img.shields.io/badge/gitee-reference-red?logo=gitee&logoColor=red&labelColor=white)](https://gitee.com/leegene/gession)
[![Github link address](https://img.shields.io/badge/github-reference-blue?logo=github&logoColor=black&labelColor=white&color=black)](https://github.com/wangluozhe/gession)
[![Go Version](https://img.shields.io/badge/Go%20Version-1.15.6-blue?logo=go&logoColor=white&labelColor=gray)]()
[![Release Version](https://img.shields.io/badge/release-v1.0.0-blue)]()

1. Gession is a conversation framework.
2. Session persistence that can be used for HTTP / HTTPS network connections.
3. The built-in docking redis interface connects and saves sessions to redis.
4. The operation is extremely convenient, easy to use, the source code is easy to understand and annotated.

# Installation
```shell script
$ go get -u github.com/wangluozhe/gession
```

# Dependent
[![Redigo v1.8.5](https://img.shields.io/badge/redigo-v1.8.5-blue?logo=go&logoColor=white)](https://github.com/gomodule/redigo)
[![go.uuid v1.2.0](https://img.shields.io/badge/go.uuid-v1.2.0-blue?logo=go&logoColor=white)](https://github.com/satori/go.uuid)

# Import
```go
import(
    "github.com/wangluozhe/gession/session"
)
```

# Init Session Manager
The session manager is used to manage various session operations, such as new / get / del.
Session manager variable name: `session.Ssmgr`.

`Redis method：`

```go
host := "host"				    // Redis address
port := 6379				    // Redis port
password := "password"		    // Redis password
database := 1				    // Redis database
pool := session.NewRedisPool(host, int(port), password, int(database))		// Create redis connect pool
expire := 1000				    // Session expire time
session.Init(pool, expire)	    // Init session manager
Ssmgr := session.Ssmgr		    // Available globally. session.Ssmgr can be operated directly or after assignment. Ssmgr is the session group manager used to create/read/delete sessions
```

`Memory method：`
```go
session.Init()                  // Init session manager
Ssmgr := session.Ssmgr		    // Available globally. session.Ssmgr can be operated directly or after assignment. Ssmgr is the session group manager used to create/read/delete sessions
```

### New Session

```go
ss := Ssmgr.New()			    // New session
```

### Get Session
```go
ss,err := Ssmgr.Get(sessionId)  // Read the session. If there is no session in memory, it will be read from redis without returning nil
if err != nil{
    fmt.Println(err)
}
```

### Delete Session
```go
ss,err := Ssmgr.Del(sessionId)  // When a session is deleted, the session in memory and redis will be deleted
if err != nil{
    fmt.Println(err)
}
```

# Session Operation
There are only four types of session operations, which are simple and convenient. They store map types: set/get/del/save.

### Set Session Value
```go
ss.Set(key,value)               // Set session value
// 或
isSuccess := ss.Set(key,value)  // Returns a bool value indicating whether the setting was successful
```

### Get Session Value
```go
result := ss.Get(key)           // Get session value
```

### Delete Session Value
```go
ss.Del(key)                     // Delete session value
// 或
isSuccess := ss.Del(key)        // Returns a bool value indicating whether the deletion was successful
```

### Save Session To Redis
```go
isSuccess,err := ss.Save()      // Save session to redis
if err != nil{
    fmt.Println(err)
}
```