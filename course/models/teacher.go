package models

import (
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Age     int    `json:"age"`
	Subject string `json:"subject"`

	IDcard string `json:"idcard" gorm:"unique"`
	Phone  string `json:"phone"`
}

func (t *Teacher) CreateTeacher() error {
	result := db.Create(t)
	return result.Error
}

func (t *Teacher) UpdateTeacher() (*Teacher, error) {
	var teacher Teacher
	result := db.Find(&teacher, t.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, nil
	}

	result = db.Save(t)
	if result.Error != nil {
		return nil, result.Error
	}

	return &teacher, nil
}

func DeleteTeacherById(Id uint) (*Teacher, error) {
	var teacher Teacher
	result := db.Delete(&teacher, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, nil
	}
	return &teacher, nil
}

func GetAllTeachers() ([]Teacher, error) {
	var teachers []Teacher
	result := db.Find(&teachers)
	return teachers, result.Error
}

func GetTeacherById(Id uint) (*Teacher, error) {
	var teacher Teacher
	result := db.Find(&teacher, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, nil
	}
	return &teacher, nil
}
