package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/youngshawn/go-project-demo/course/config"
	"github.com/youngshawn/go-project-demo/course/models"
)

func initGinContext(ctx *gin.Context) {
	// Dynamic config
	//ctx.Set("Config.Cache.EnableNullResultCache", config.EnableNullResultCache)
}

func Status(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(http.StatusOK)
	status, details := models.Healthcheck()
	json.NewEncoder(ctx.Writer).Encode(map[string]interface{}{
		"status":  status,
		"details": details,
	})
}

func Version(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(map[string]string{"version": config.Version})
}
