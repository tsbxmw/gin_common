package middleware

import (
    yaag_gin "github.com/betacraft/yaag/gin"
    "github.com/betacraft/yaag/yaag"
    "github.com/gin-gonic/gin"
)

func ApidocMiddlewareInit(e *gin.Engine) {
    yaag.Init(&yaag.Config{On: true, DocTitle: "Gin", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})
    apidocMiddleware := yaag_gin.Document()
    e.Use(apidocMiddleware)
}
