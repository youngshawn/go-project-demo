package controllers

import (
	"github.com/gin-gonic/gin"
)

func initGinContext(ctx *gin.Context) {
	// Dynamic config
	//ctx.Set("Config.Cache.EnableNullResultCache", config.EnableNullResultCache)
}
