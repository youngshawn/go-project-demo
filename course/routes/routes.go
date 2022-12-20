package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/youngshawn/go-project-demo/course/controllers"
)

func InstallRoutes(router *gin.Engine) {
	router.GET("/status", controllers.Status)
	router.GET("/version", controllers.Version)
	installCourseRouters(router)
	installTeacherRouters(router)
}
