package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/youngshawn/go-project-demo/course/controllers"
)

func installTeacherRouters(router *gin.Engine) {
	g := router.Group("/teacher")
	{
		g.GET("/", controllers.GetAllTeachers)
		g.POST("/", controllers.CreateTeacher)
		g.GET("/:id", controllers.GetTeacherById)
		g.PUT("/:id", controllers.UpdateTeacherById)
		g.DELETE("/:id", controllers.DeleteTeacherById)

		g.GET("/:id/course", controllers.GetCoursesByTeacherId)
	}
}
