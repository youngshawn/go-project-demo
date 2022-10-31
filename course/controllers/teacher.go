package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/youngshawn/go-project-demo/course/models"
)

func CreateTeacher(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	teacher := &models.Teacher{}
	json.NewDecoder(ctx.Request.Body).Decode(teacher)
	err := teacher.CreateTeacher()
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to create teacher")
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(teacher)
}

func UpdateTeacherById(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	//id := ctx.GetUint("id")
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 0, 64)
	if err != nil {
		log.Printf("The parameter <id> is mailformat: %s\n", idParam)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode("Error: teacher id is malformat")
		return
	}
	teacher := &models.Teacher{}
	json.NewDecoder(ctx.Request.Body).Decode(teacher)
	teacher.ID = uint(id)
	_, err = teacher.UpdateTeacher()
	if err == models.ErrorObjectNotFound {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(ctx.Writer).Encode("Error: teacher not found")
		return
	}
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to update teacher")
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode("Teacher updated")
}

func DeleteTeacherById(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	//id := ctx.GetUint("id")
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 0, 64)
	if err != nil {
		log.Printf("The parameter <id> is mailformat: %s\n", idParam)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode("Error: Teacher id is malformat")
		return
	}
	_, err = models.DeleteTeacherById(uint(id))
	if err == models.ErrorObjectNotFound {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(ctx.Writer).Encode("Error: Teacher not found")
		return
	}
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to delete teacher")
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode("Teacher deleted")
}

func GetAllTeachers(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	teachers, err := models.GetAllTeachers()
	if err == models.ErrorObjectNotFound {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		//json.NewEncoder(ctx.Writer).Encode("Error: Teachers not found")
		json.NewEncoder(ctx.Writer).Encode(teachers)
		return
	}
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to get all teachers")
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(teachers)
}

func GetTeacherById(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	//id := ctx.GetUint("id")
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 0, 64)
	if err != nil {
		log.Printf("The parameter <id> is mailformat: %s\n", idParam)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ctx.Writer).Encode("Error: Teacher id is malformat")
		return
	}
	teacher, err := models.GetTeacherById(uint(id))
	if err == models.ErrorObjectNotFound {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		//json.NewEncoder(ctx.Writer).Encode("Error: Teacher not found")
		json.NewEncoder(ctx.Writer).Encode(nil)
		return
	}
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to get teacher")
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(teacher)
}
