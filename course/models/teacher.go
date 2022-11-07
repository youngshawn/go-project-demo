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
	Cache.DeleteFromLocalCache("/teacher/")
	Cache.Delete(context.Background(), "/teacher/")

	result := db.Create(t)

	time.Sleep(time.Millisecond * 500)
	Cache.DeleteFromLocalCache("/teacher/")
	Cache.Delete(context.Background(), "/teacher/")
	Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", t.ID))
	Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", t.ID))
	Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", t.ID))
	Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", t.ID))

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

	Cache.DeleteFromLocalCache("/teacher/")
	Cache.Delete(context.Background(), "/teacher/")
	Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", t.ID))
	Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", t.ID))
	courses, _ := GetCoursesByTeacherId(t.ID)
	for _, c := range courses {
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
	}

	result = db.Save(t)

	time.Sleep(time.Millisecond * 500)
	Cache.DeleteFromLocalCache("/teacher/")
	Cache.Delete(context.Background(), "/teacher/")
	Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", t.ID))
	Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", t.ID))
	for _, c := range courses {
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
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

	Cache.DeleteFromLocalCache("/teacher/")
	Cache.Delete(context.Background(), "/teacher/")
	Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", Id))
	Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", Id))
	Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", Id))
	Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", Id))
	courses, _ := GetCoursesByTeacherId(Id)
	for _, c := range courses {
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
	}

	result = db.Delete(&teacher, Id)

	time.Sleep(time.Millisecond * 500)
	Cache.DeleteFromLocalCache("/teacher/")
	Cache.Delete(context.Background(), "/teacher/")
	Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", Id))
	Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", Id))
	Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", Id))
	Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", Id))
	for _, c := range courses {
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
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
	dynamicCacheConfig := config.GetDynamicCacheConfig()

	var teachers []Teacher

	err := Cache.Once(&cache.Item{
		Key:            "/teacher/",
		Value:          &teachers,
		TTL:            dynamicCacheConfig.CacheTTL,
		SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Println("GetAllTeachers from DB...")
			result := db.Find(i.Value)

			err := result.Error
			if err == nil && result.RowsAffected == 0 && !dynamicCacheConfig.EnableNullResultCache {
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
	dynamicCacheConfig := config.GetDynamicCacheConfig()

	var teacher Teacher

	err := Cache.Once(&cache.Item{
		Key:            fmt.Sprintf("/teacher/%d", Id),
		Value:          &teacher,
		TTL:            dynamicCacheConfig.CacheTTL,
		SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Println("GetTeacherById from DB...")
			result := db.Find(i.Value, Id)

			err := result.Error
			if err == nil && result.RowsAffected == 0 { // Null Result
				if dynamicCacheConfig.EnableNullResultCache {
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
