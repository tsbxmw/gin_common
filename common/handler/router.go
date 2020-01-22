package handler

import (
	"github.com/gin-gonic/gin"
)

func NoRouteHandlerInit(e *gin.Engine) {
	noRouteHandler := NoRouteHandler()
	e.NoRoute(noRouteHandler)
}

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "page not found",
			"data":    []string{},
		})
	}
}
