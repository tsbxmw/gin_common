package handler

import (
	"github.com/gin-gonic/gin"
)

func NoMethodHandlerInit(e *gin.Engine) {
	noMethodHandler := NoMethodHandler()
	e.NoMethod(noMethodHandler)
}

func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(405, gin.H{
			"code":    405,
			"message": "method not allowed",
			"data":    []string{},
		})
	}
}
