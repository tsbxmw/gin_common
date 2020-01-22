package middleware

import (
    "github.com/gin-gonic/gin"
    common "github.com/tsbxmw/gin_common"
)

func HeaderMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        common.InitKey(c)
        tcsdkVersion := c.GetHeader("tcsdk_version")
        c.Keys["tcsdk_version"] = tcsdkVersion
        c.Next()
    }
}

// global auth middleware init
func HeaderMiddlewareInit(e *gin.Engine) {
    header := HeaderMiddleware()
    e.Use(header)
}
