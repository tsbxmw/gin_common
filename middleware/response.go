package middleware

import (
    "github.com/gin-gonic/gin"
    common "github.com/tsbxmw/gin_common"
)

func ResponseMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if c.Writer.Written() {
            return
        }
        c.JSON(common.HTTP_STATUS_OK, common.Response{
            Code: 0,
            Message: "success",
            Data: []string{},
        })
    }
}

// global response middleware init
func ResponseMiddlewareInit(e *gin.Engine) {
    response := ResponseMiddleware()
    e.Use(response)
}
