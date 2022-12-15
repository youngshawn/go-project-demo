package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/youngshawn/go-project-demo/course/controllers"
)

func InstallRoutes(router *gin.Engine) {
	router.GET("/health", controllers.Health)
	installCourseRouters(router)
	installTeacherRouters(router)
}
