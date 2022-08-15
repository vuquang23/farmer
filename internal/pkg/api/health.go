package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthApi struct {
}

func newHealthApi() *healthApi {
	return &healthApi{}
}

func (a *healthApi) setupRoute(rg *gin.RouterGroup) {
	rg.GET("/live", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, "OK")
	})
}
