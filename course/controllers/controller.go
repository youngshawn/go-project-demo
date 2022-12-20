package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youngshawn/go-project-demo/course/config"
)

func initGinContext(ctx *gin.Context) {
	// Dynamic config
	//ctx.Set("Config.Cache.EnableNullResultCache", config.EnableNullResultCache)
}

func Status(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(map[string]string{"status": "healthy"})
}

func Version(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(map[string]string{"version": config.Version})
}
