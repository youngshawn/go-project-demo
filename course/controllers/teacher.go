package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/youngshawn/go-project-demo/course/models"
)

func GetAllTeachers(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json")

	teachers, err := models.GetAllTeachers()
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
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to get teacher")
		return
	}
	if teacher == nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(ctx.Writer).Encode("Error: Teacher not found")
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode(teacher)
}

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
		json.NewEncoder(ctx.Writer).Encode("Error: Teacher id is malformat")
		return
	}
	teacher := &models.Teacher{}
	json.NewDecoder(ctx.Request.Body).Decode(teacher)
	teacher.ID = uint(id)
	teacher, err = teacher.UpdateTeacher()
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to update teacher")
		return
	}

	if teacher == nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(ctx.Writer).Encode("Error: Teacher not found")
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
	teacher, err := models.DeleteTeacherById(uint(id))
	if err != nil {
		log.Println(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ctx.Writer).Encode("Error: Failed to delete teacher")
		return
	}

	if teacher == nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		json.NewEncoder(ctx.Writer).Encode("Error: Teacher not found")
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(ctx.Writer).Encode("Teacher deleted")
}
