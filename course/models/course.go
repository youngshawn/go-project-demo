package models

import (
	"gorm.io/gorm"
)

// Course belongs to Teacher
type Course struct {
	gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"author"`
	Hours       uint    `json:"hours"`
	Room        string  `json:"room"`
	Duration    string  `json:"duration"`
	Teacher     Teacher `json:"teacher" gorm:"foreignKey:TeacherID"`
	TeacherID   uint
}

func GetAllCourses() []Course {
	var courses []Course
	db.Find(&courses)
	return courses
}
