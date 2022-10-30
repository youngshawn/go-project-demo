package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/youngshawn/go-project-demo/course/models"
)

func CreateCourse(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	course := &models.Course{}
	json.NewDecoder(ctx.Request.Body).Decode(course)
	err := course.CreateCourse()
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to create course")
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(course)
}

func UpdateCourseById(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	//id := ctx.GetUint("id")
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 0, 64)
	if err != nil {
		log.Printf("The parameter <id> is mailformat: %s\n", idParam)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode("Error: course id is malformat")
		return
	}
	course := &models.Course{}
	json.NewDecoder(ctx.Request.Body).Decode(course)
	course.ID = uint(id)
	_, err = course.UpdateCourse()
	if err == models.ErrorObjectNotFound {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(ctx.Writer).Encode("Error: Course not found")
		return
	}
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to update course")
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode("Course updated")
}

func DeleteCourseById(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	//id := ctx.GetUint("id")
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 0, 64)
	if err != nil {
		log.Printf("The parameter <id> is mailformat: %s\n", idParam)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode("Error: Course id is malformat")
		return
	}
	_, err = models.DeleteCourseById(uint(id))
	if err == models.ErrorObjectNotFound {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(ctx.Writer).Encode("Error: Course not found")
		return
	}
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to delete course")
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode("Course deleted")
}

func GetAllCourses(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	courses, err := models.GetAllCourses()
	if err == models.ErrorObjectNotFound {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		//json.NewEncoder(ctx.Writer).Encode("Error: Courses not found")
		json.NewEncoder(ctx.Writer).Encode(courses)
		return
	}
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to get all courses")
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(courses)
}

func GetCourseById(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	//id := ctx.GetUint("id")
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 0, 64)
	if err != nil {
		log.Printf("The parameter <id> is mailformat: %s\n", idParam)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode("Error: Course id is malformat")
		return
	}
	course, err := models.GetCourseById(uint(id))
	if err == models.ErrorObjectNotFound {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		//json.NewEncoder(ctx.Writer).Encode("Error: Course not found")
		json.NewEncoder(ctx.Writer).Encode(nil)
		return
	}
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to get course")
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(course)
}

func GetCoursesByTeacherId(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	teacherIdParam := ctx.Param("id")
	teacherId, err := strconv.ParseUint(teacherIdParam, 0, 64)
	if err != nil {
		log.Printf("The parameter <id> is mailformat: %s\n", teacherIdParam)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode("Error: Teacher id is malformat")
		return
	}
	courses, err := models.GetCoursesByTeacherId(uint(teacherId))
	if err == models.ErrorObjectNotFound {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		//json.NewEncoder(ctx.Writer).Encode("Error: Courses not found")
		json.NewEncoder(ctx.Writer).Encode(courses)
		return
	}
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to get courses by teacher")
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(courses)
}

func GetTeacherByCourseId(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	courseIdParam := ctx.Param("id")
	courseId, err := strconv.ParseUint(courseIdParam, 0, 64)
	if err != nil {
		log.Printf("The parameter <id> is mailformat: %s\n", courseIdParam)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode("Error: Course id is malformat")
		return
	}
	teacher, err := models.GetTeacherByCourseId(uint(courseId))
	if err == models.ErrorObjectNotFound {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		//json.NewEncoder(ctx.Writer).Encode("Error: Teacher not found")
		json.NewEncoder(ctx.Writer).Encode(nil)
		return
	}
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to get teacher by course")
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(teacher)
}
