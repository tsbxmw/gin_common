package handler

import "github.com/gin-gonic/gin"

func HandlerInit(engin *gin.Engine) {
	NoRouteHandlerInit(engin)
	NoMethodHandlerInit(engin)
}
