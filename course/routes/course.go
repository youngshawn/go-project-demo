package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/youngshawn/go-project-demo/course/controllers"
)

func installCourseRouters(router *gin.Engine) {
	course := router.Group("/course")
	{
		course.GET("/", controllers.GetAllCourses)
		course.POST("/", controllers.CreateCourse)
		course.GET("/:id", controllers.GetCourseById)
		course.PUT("/:id", controllers.UpdateCourseById)
		course.DELETE("/:id", controllers.DeleteCourseById)

		course.GET("/:id/teacher", controllers.GetTeacherByCourseId)
	}
}
