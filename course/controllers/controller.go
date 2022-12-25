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

	status, details := models.Healthcheck()
	if status == "healthy" {
		ctx.Writer.WriteHeader(http.StatusOK)
	} else {
		ctx.Writer.WriteHeader(http.StatusServiceUnavailable)
	}

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
