package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initGinContext(ctx *gin.Context) {
	// Dynamic config
	//ctx.Set("Config.Cache.EnableNullResultCache", config.EnableNullResultCache)
}

func Health(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(map[string]string{"status": "healthy"})
}
