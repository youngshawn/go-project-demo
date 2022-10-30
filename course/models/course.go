package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/youngshawn/go-project-demo/course/config"
	"gorm.io/gorm"
)

var ErrorObjectNotFound = errors.New("ObjectNotFound")

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
	redisCache.DeleteFromLocalCache("/course/")
	redisCache.Delete(context.Background(), "/course/")

	result := db.Create(c)

	time.Sleep(time.Millisecond * 500)
	redisCache.DeleteFromLocalCache("/course/")
	redisCache.Delete(context.Background(), "/course/")
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d", c.ID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d", c.ID))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", c.TeacherID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", c.TeacherID))

	return result.Error
}

func (c *Course) UpdateCourse() (*Course, error) {
	var course Course

	result := db.Find(&course, c.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrorObjectNotFound
	}

	redisCache.DeleteFromLocalCache("/course/")
	redisCache.Delete(context.Background(), "/course/")
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d", c.ID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d", c.ID))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", c.TeacherID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", c.TeacherID))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", course.TeacherID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", course.TeacherID))

	result = db.Save(c)

	time.Sleep(time.Millisecond * 500)
	redisCache.DeleteFromLocalCache("/course/")
	redisCache.Delete(context.Background(), "/course/")
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d", c.ID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d", c.ID))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", c.TeacherID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", c.TeacherID))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", course.TeacherID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", course.TeacherID))

	return c, result.Error
}

func DeleteCourseById(Id uint) (*Course, error) {
	var course Course

	result := db.Find(&course, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrorObjectNotFound
	}

	redisCache.DeleteFromLocalCache("/course/")
	redisCache.Delete(context.Background(), "/course/")
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d", Id))
	redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d", Id))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", Id))
	redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", Id))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", course.TeacherID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", course.TeacherID))

	result = db.Delete(&course, Id)

	time.Sleep(time.Millisecond * 500)
	redisCache.DeleteFromLocalCache("/course/")
	redisCache.Delete(context.Background(), "/course/")
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d", Id))
	redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d", Id))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", Id))
	redisCache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", Id))
	redisCache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", course.TeacherID))
	redisCache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", course.TeacherID))

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrorObjectNotFound
	}
	return &course, nil
}

func GetAllCourses() ([]Course, error) {
	var courses []Course

	err := redisCache.Once(&cache.Item{
		Key:            "/course/",
		Value:          &courses,
		TTL:            config.CacheTTL,
		SkipLocalCache: !config.EnableLocalCache,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Println("GetAllCourses from DB...")
			result := db.Find(i.Value)

			err := result.Error
			if err == nil && result.RowsAffected == 0 && !config.EnableNullResultCache {
				err = ErrorObjectNotFound
			}

			return i.Value, err
		},
	})

	if err == nil && len(courses) == 0 {
		err = ErrorObjectNotFound
	}

	return courses, err
}

func GetCourseById(Id uint) (*Course, error) {
	var course Course

	err := redisCache.Once(&cache.Item{
		Key:            fmt.Sprintf("/course/%d", Id),
		Value:          &course,
		TTL:            config.CacheTTL,
		SkipLocalCache: !config.EnableLocalCache,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Println("GetCourseById from DB...")
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

	if err == ErrorObjectNotFound || err == nil && course.ID == 0 {
		return nil, ErrorObjectNotFound
	}

	return &course, err
}

func GetCoursesByTeacherId(Id uint) ([]Course, error) {
	var courses []Course

	err := redisCache.Once(&cache.Item{
		Key:            fmt.Sprintf("/teacher/%d/course", Id),
		Value:          &courses,
		TTL:            config.CacheTTL,
		SkipLocalCache: !config.EnableLocalCache,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Println("GetCoursesByTeacherId from DB...")

			result := db.Where(&Course{TeacherID: Id}).Find(i.Value)

			err := result.Error
			if err == nil && result.RowsAffected == 0 && !config.EnableNullResultCache {
				err = ErrorObjectNotFound
			}

			return i.Value, err
		},
	})

	if err == nil && len(courses) == 0 {
		err = ErrorObjectNotFound
	}

	return courses, err
}

func GetTeacherByCourseId(Id uint) (*Teacher, error) {
	var teacher Teacher
	err := redisCache.Once(&cache.Item{
		Key:            fmt.Sprintf("/course/%d/teacher", Id),
		Value:          &teacher,
		TTL:            config.CacheTTL,
		SkipLocalCache: !config.EnableLocalCache,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Println("GetTeacherByCourseId from DB...")
			var course Course
			err1 := redisCache.Once(&cache.Item{
				Key:            fmt.Sprintf("/course/%d", Id),
				Value:          &course,
				TTL:            config.CacheTTL,
				SkipLocalCache: !config.EnableLocalCache,
				Do: func(i *cache.Item) (interface{}, error) {
					log.Println("GetCourseById from DB...")

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

			if err1 == ErrorObjectNotFound || err1 == nil && (course.ID == 0 || course.TeacherID == 0) { // Null Result
				if config.EnableNullResultCache {
					return nil, nil
				} else {
					return nil, ErrorObjectNotFound
				}
			}
			if err1 != nil {
				return nil, err1
			}

			var teacher Teacher
			err2 := redisCache.Once(&cache.Item{
				Key:            fmt.Sprintf("/teacher/%d", course.TeacherID),
				Value:          &teacher,
				TTL:            config.CacheTTL,
				SkipLocalCache: !config.EnableLocalCache,
				Do: func(i *cache.Item) (interface{}, error) {
					log.Println("GetTeacherById from DB...")
					var teacher Teacher
					err3 := db.Model(&course).Association("Teacher").Find(&teacher)

					if err3 == nil && teacher.ID == 0 { // Null Result
						if config.EnableNullResultCache {
							return nil, nil
						} else {
							return nil, ErrorObjectNotFound
						}
					}

					return &teacher, nil
				},
			})

			if err2 == ErrorObjectNotFound || err2 == nil && teacher.ID == 0 { // Null Result
				if config.EnableNullResultCache {
					return nil, nil
				} else {
					return nil, ErrorObjectNotFound
				}
			}

			return &teacher, err2
		},
	})

	if err == ErrorObjectNotFound || err == nil && teacher.ID == 0 {
		return nil, ErrorObjectNotFound
	}

	return &teacher, err
}
