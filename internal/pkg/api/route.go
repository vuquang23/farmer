package api

import (
	"github.com/gin-gonic/gin"
)

type IAPI interface {
	setupRoute(rg *gin.RouterGroup)
}

func AddRouterV1(server *gin.Engine) {
	// API v1
	v1 := server.Group("/api/v1")

	addApi(newHealthApi(), "/health", v1)
}

func addApi(api IAPI, path string, rg *gin.RouterGroup) {
	apiRg := rg.Group(path)
	api.setupRoute(apiRg)
}
