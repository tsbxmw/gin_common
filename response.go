package common

import "github.com/gin-gonic/gin"

type (
    Response struct {
        Code    int         `json:"code"`
        Message string      `json:"msg"`
        Data    interface{} `json:"data"`
    }
    GinResponse struct {
        Ctx *gin.Context
    }
)

func (gr *GinResponse) Response(code int, msg string, data interface{}) {
    gr.Ctx.JSON(200, Response{
        Code:    code,
        Message: msg,
        Data:    data,
    })
    return
}
