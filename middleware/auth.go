package middleware

import (
    common "github.com/tsbxmw/gin_common"
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authGlobal := common.AuthGlobal{}
        common.InitKey(c)
        //这一部分可以替换成从session/cookie中获取，
        key := c.GetHeader("key")
        secret := c.GetHeader("secret")
        //var userId int
        if key == "" {
            key = c.Query("key")
        }
        if secret == "" {
            secret = c.Query("secret")
        }

        if key == "" || secret == "" {
            c.Keys["code"] = common.HTTP_AUTH_ERROR
            panic(common.NewHttpAuthError())
        }

        // redis get first
        redisConn := common.RedisPool.Get()
        defer redisConn.Close()

        authRedis := common.AuthRedis{}
        auth := common.AuthModel{}

        redisFlag := false

        if secretTemp, err := common.RedisGet(c, redisConn, key); err != nil {
            common.LogrusLogger.Error(err)
            if err := common.DB.Table("auth").Where("app_key=? and app_secret=?", key, secret).First(&auth).Error; err != nil {
                c.Keys["code"] = common.REDIS_GET_ERROR
                panic(err)
            }
        } else {
            fmt.Println(secretTemp)
            if err = json.Unmarshal([]byte(secretTemp), &authRedis); err != nil {
                c.Keys["code"] = common.REDIS_GET_ERROR
                panic(err)
            }
            if secret != authRedis.Secret {
                c.Keys["code"] = common.HTTP_AUTH_ERROR
                panic(common.NewHttpAuthError())
            }
            redisFlag = true
        }
        if redisFlag {
            authGlobal.UserId = authRedis.UserId
        } else {
            authGlobal.UserId = auth.UserId
        }
        c.Keys["auth"] = &authGlobal
        c.Next()
    }
}

// global auth middleware init
func AuthInit(e *gin.Engine) {
    auth := AuthMiddleware()
    e.Use(auth)
}
