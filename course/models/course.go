package models

import (
	"gorm.io/gorm"
)

// Course belongs to Teacher
type Course struct {
	gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Hours       uint    `json:"hours"`
	Room        string  `json:"room"`
	Duration    string  `json:"duration"`
	TeacherID   uint    `json:"teacher"`
	Teacher     Teacher `json:"-" gorm:"foreignKey:TeacherID"`
}

func (c *Course) CreateCourse() error {
	result := db.Create(c)
	return result.Error
}

func (c *Course) UpdateCourse() (*Course, error) {
	var course Course
	result := db.Find(&course, c.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, nil
	}

	result = db.Save(c)
	if result.Error != nil {
		return nil, result.Error
	}

	return &course, nil
}

func GetAllCourses() ([]Course, error) {
	var courses []Course
	result := db.Find(&courses)
	return courses, result.Error
}

func GetCoursesByTeacherId(Id uint) ([]Course, error) {
	var courses []Course
	result := db.Where(&Course{TeacherID: Id}).Find(&courses)
	return courses, result.Error
}

func GetCourseById(Id uint) (*Course, error) {
	var course Course
	result := db.Find(&course, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, nil
	}
	return &course, nil
}

func DeleteCourseById(Id uint) (*Course, error) {
	var course Course
	result := db.Delete(&course, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, nil
	}
	return &course, nil
}
