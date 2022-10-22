package routes

import "github.com/gin-gonic/gin"

func InstallRoutes(router *gin.Engine) {
	installCourseRouters(router)
	installTeacherRouters(router)
}
