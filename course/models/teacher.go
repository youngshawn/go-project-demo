package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/youngshawn/go-project-demo/course/config"
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
	redisCache.DeleteFromLocalCache("/teacher/")
	redisCache.Delete(context.Background(), "/teacher/")

	result := db.Create(t)

	time.Sleep(time.Millisecond * 500)
	redisCache.DeleteFromLocalCache("/teacher/")
	redisCache.Delete(context.Background(), "/teacher/")
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", t.ID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", t.ID))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", t.ID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", t.ID))

	return result.Error
}

func (t *Teacher) UpdateTeacher() (*Teacher, error) {
	var teacher Teacher

	result := db.Find(&teacher, t.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrorObjectNotFound
	}

	redisCache.DeleteFromLocalCache("/teacher/")
	redisCache.Delete(context.Background(), "/teacher/")
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", t.ID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", t.ID))
	courses, _ := GetCoursesByTeacherId(t.ID)
	for _, c := range courses {
		redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
	}

	result = db.Save(t)

	time.Sleep(time.Millisecond * 500)
	redisCache.DeleteFromLocalCache("/teacher/")
	redisCache.Delete(context.Background(), "/teacher/")
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", t.ID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", t.ID))
	for _, c := range courses {
		redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
	}

	return t, result.Error
}

func DeleteTeacherById(Id uint) (*Teacher, error) {
	var teacher Teacher

	result := db.Find(&teacher, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrorObjectNotFound
	}

	redisCache.DeleteFromLocalCache("/teacher/")
	redisCache.Delete(context.Background(), "/teacher/")
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", Id))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", Id))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", Id))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", Id))
	courses, _ := GetCoursesByTeacherId(Id)
	for _, c := range courses {
		redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
	}

	result = db.Delete(&teacher, Id)

	time.Sleep(time.Millisecond * 500)
	redisCache.DeleteFromLocalCache("/teacher/")
	redisCache.Delete(context.Background(), "/teacher/")
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", Id))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", Id))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", Id))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", Id))
	for _, c := range courses {
		redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
	}

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrorObjectNotFound
	}
	return &teacher, nil
}

func GetAllTeachers() ([]Teacher, error) {
	var teachers []Teacher

	err := redisCache.Once(&cache.Item{
		Key:            "/teacher/",
		Value:          &teachers,
		TTL:            config.CacheTTL,
		SkipLocalCache: !config.EnableLocalCache,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Println("GetAllTeachers from DB...")
			result := db.Find(i.Value)

			err := result.Error
			if err == nil && result.RowsAffected == 0 && !config.EnableNullResultCache {
				err = ErrorObjectNotFound
			}

			return i.Value, err
		},
	})

	if err == nil && len(teachers) == 0 {
		err = ErrorObjectNotFound
	}

	return teachers, err
}

func GetTeacherById(Id uint) (*Teacher, error) {
	var teacher Teacher

	err := redisCache.Once(&cache.Item{
		Key:            fmt.Sprintf("/teacher/%d", Id),
		Value:          &teacher,
		TTL:            config.CacheTTL,
		SkipLocalCache: !config.EnableLocalCache,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Println("GetTeacherById from DB...")
			result := db.Find(i.Value, Id)

			err := result.Error
			if err == nil && result.RowsAffected == 0 { // Null Result
				if config.EnableNullResultCache {
					return nil, nil
				} else {
					return nil, ErrorObjectNotFound
				}
			}
			return i.Value, err
		},
	})

	if err == ErrorObjectNotFound || err == nil && teacher.ID == 0 {
		return nil, ErrorObjectNotFound
	}

	return &teacher, err
}
