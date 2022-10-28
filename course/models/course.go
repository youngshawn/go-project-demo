package models

import (
	"github.com/go-redis/cache/v9"
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

	err := redisCache.Once(&cache.Item{
		Key:            "/course/",
		Value:          &courses,
		SkipLocalCache: false,
		Do: func(i *cache.Item) (interface{}, error) {
			//var cs []Course
			//result := db.Find(&cs)
			//return &cs, result.Error
			result := db.Find(i.Value)
			return i.Value, result.Error
		},
	})

	return courses, err
}

func GetCoursesByTeacherId(Id uint) ([]Course, error) {
	var courses []Course
	result := db.Where(&Course{TeacherID: Id}).Find(&courses)
	return courses, result.Error
}

func GetTeacherByCourseId(Id uint) (*Teacher, error) {
	var course Course
	result := db.Find(&course, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected != 1 {
		return nil, nil
	}
	var teacher Teacher
	err := db.Model(&course).Association("Teacher").Find(&teacher)
	if err != nil {
		return nil, err
	}
	if teacher.ID == 0 {
		return nil, nil
	}
	return &teacher, nil
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
