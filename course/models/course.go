package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-redis/cache/v9"
	"github.com/youngshawn/go-project-demo/course/config"
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

	err := hystrix.Do("CourseCUD", func() error {

		Cache.DeleteFromLocalCache("/course/")
		Cache.Delete(context.Background(), "/course/")

		result := db.Create(c)

		time.Sleep(time.Millisecond * 500)
		Cache.DeleteFromLocalCache("/course/")
		Cache.Delete(context.Background(), "/course/")
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d", c.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d", c.ID))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", c.TeacherID))
		Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", c.TeacherID))

		return result.Error

	}, nil)

	return err
}

func (c *Course) UpdateCourse() (*Course, error) {
	var course Course

	err := hystrix.Do("CourseCUD", func() error {

		result := db.Find(&course, c.ID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrorObjectNotFound
		}

		Cache.DeleteFromLocalCache("/course/")
		Cache.Delete(context.Background(), "/course/")
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d", c.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d", c.ID))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", c.TeacherID))
		Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", c.TeacherID))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", course.TeacherID))
		Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", course.TeacherID))

		result = db.Omit("CreatedAt").Save(c)

		time.Sleep(time.Millisecond * 500)
		Cache.DeleteFromLocalCache("/course/")
		Cache.Delete(context.Background(), "/course/")
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d", c.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d", c.ID))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", c.TeacherID))
		Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", c.TeacherID))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", course.TeacherID))
		Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", course.TeacherID))

		return result.Error

	}, nil)

	if err != nil {
		return nil, err
	}
	return c, nil
}

func DeleteCourseById(Id uint) (*Course, error) {
	var course Course

	err := hystrix.Do("CourseCUD", func() error {

		result := db.Find(&course, Id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrorObjectNotFound
		}

		Cache.DeleteFromLocalCache("/course/")
		Cache.Delete(context.Background(), "/course/")
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d", Id))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d", Id))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", Id))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", Id))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", course.TeacherID))
		Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", course.TeacherID))

		result = db.Delete(&course, Id)

		time.Sleep(time.Millisecond * 500)
		Cache.DeleteFromLocalCache("/course/")
		Cache.Delete(context.Background(), "/course/")
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d", Id))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d", Id))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", Id))
		Cache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", Id))
		Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d/course", course.TeacherID))
		Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d/course", course.TeacherID))

		return result.Error

	}, nil)

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func GetAllCourses() ([]Course, error) {
	dynamicCacheConfig := config.GetDynamicCacheConfig()

	var courses []Course

	err := hystrix.Do("GetAllCourses", func() error {
		return Cache.Once(&cache.Item{
			Key:            "/course/",
			Value:          &courses,
			TTL:            dynamicCacheConfig.CacheTTL,
			SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
			Do: func(i *cache.Item) (interface{}, error) {
				log.Println("GetAllCourses: miss cache get from DB...")
				result := db.Find(i.Value)

				err := result.Error
				if err == nil && result.RowsAffected == 0 && !dynamicCacheConfig.EnableNullResultCache {
					err = ErrorObjectNotFound
				}

				return i.Value, err
			},
		})
	}, nil)

	if err == nil && len(courses) == 0 {
		err = ErrorObjectNotFound
	}

	return courses, err
}

func GetCourseById(Id uint) (*Course, error) {
	dynamicCacheConfig := config.GetDynamicCacheConfig()

	var course Course

	err := hystrix.Do("GetCourseById", func() error {
		return Cache.Once(&cache.Item{
			Key:            fmt.Sprintf("/course/%d", Id),
			Value:          &course,
			TTL:            dynamicCacheConfig.CacheTTL,
			SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
			Do: func(i *cache.Item) (interface{}, error) {
				log.Println("GetCourseById: miss cache get from DB...")
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
	}, nil)

	if err == ErrorObjectNotFound || err == nil && course.ID == 0 {
		return nil, ErrorObjectNotFound
	}

	return &course, err
}

func GetCoursesByTeacherId(Id uint) ([]Course, error) {
	dynamicCacheConfig := config.GetDynamicCacheConfig()

	var courses []Course

	err := hystrix.Do("GetCoursesByTeacherId", func() error {
		return Cache.Once(&cache.Item{
			Key:            fmt.Sprintf("/teacher/%d/course", Id),
			Value:          &courses,
			TTL:            dynamicCacheConfig.CacheTTL,
			SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
			Do: func(i *cache.Item) (interface{}, error) {
				log.Println("GetCoursesByTeacherId: miss cache get from DB...")

				result := db.Where(&Course{TeacherID: Id}).Find(i.Value)

				err := result.Error
				if err == nil && result.RowsAffected == 0 && !dynamicCacheConfig.EnableNullResultCache {
					err = ErrorObjectNotFound
				}

				return i.Value, err
			},
		})
	}, nil)

	if err == nil && len(courses) == 0 {
		err = ErrorObjectNotFound
	}

	return courses, err
}

func GetTeacherByCourseId2(Id uint) (*Teacher, error) {
	dynamicCacheConfig := config.GetDynamicCacheConfig()

	var teacher Teacher
	err := Cache.Once(&cache.Item{
		Key:            fmt.Sprintf("/course/%d/teacher", Id),
		Value:          &teacher,
		TTL:            dynamicCacheConfig.CacheTTL,
		SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
		Do: func(i *cache.Item) (interface{}, error) {
			log.Println("GetTeacherByCourseId: miss cache get from source...")
			var course Course
			err1 := Cache.Once(&cache.Item{
				Key:            fmt.Sprintf("/course/%d", Id),
				Value:          &course,
				TTL:            dynamicCacheConfig.CacheTTL,
				SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
				Do: func(i *cache.Item) (interface{}, error) {
					log.Println("GetCourseById: miss cache get from DB...")

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

			if err1 == ErrorObjectNotFound || err1 == nil && (course.ID == 0 || course.TeacherID == 0) { // Null Result
				if dynamicCacheConfig.EnableNullResultCache {
					return nil, nil
				} else {
					return nil, ErrorObjectNotFound
				}
			}
			if err1 != nil {
				return nil, err1
			}

			var teacher Teacher
			err2 := Cache.Once(&cache.Item{
				Key:            fmt.Sprintf("/teacher/%d", course.TeacherID),
				Value:          &teacher,
				TTL:            dynamicCacheConfig.CacheTTL,
				SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
				Do: func(i *cache.Item) (interface{}, error) {
					log.Println("GetTeacherById: miss cache get from DB...")
					var teacher Teacher
					err3 := db.Model(&course).Association("Teacher").Find(&teacher)

					if err3 == nil && teacher.ID == 0 { // Null Result
						if dynamicCacheConfig.EnableNullResultCache {
							return nil, nil
						} else {
							return nil, ErrorObjectNotFound
						}
					}

					return &teacher, nil
				},
			})

			if err2 == ErrorObjectNotFound || err2 == nil && teacher.ID == 0 { // Null Result
				if dynamicCacheConfig.EnableNullResultCache {
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

func GetTeacherByCourseId(Id uint) (*Teacher, error) {
	dynamicCacheConfig := config.GetDynamicCacheConfig()

	var teacher Teacher

	err := hystrix.Do("GetTeacherByCourseId", func() error {
		return Cache.Once(&cache.Item{
			Key:            fmt.Sprintf("/course/%d/teacher", Id),
			Value:          &teacher,
			TTL:            dynamicCacheConfig.CacheTTL,
			SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
			Do: func(i *cache.Item) (interface{}, error) {
				log.Println("GetTeacherByCourseId: miss cache get from source...")

				course, err1 := GetCourseById(Id)
				if err1 == ErrorObjectNotFound || err1 == nil && (course.ID == 0 || course.TeacherID == 0) { // Null Result
					if dynamicCacheConfig.EnableNullResultCache {
						return nil, nil
					} else {
						return nil, ErrorObjectNotFound
					}
				}
				if err1 != nil {
					return nil, err1
				}

				teacher, err2 := GetTeacherById(course.TeacherID)
				if err2 == ErrorObjectNotFound || err2 == nil && teacher.ID == 0 { // Null Result
					if dynamicCacheConfig.EnableNullResultCache {
						return nil, nil
					} else {
						return nil, ErrorObjectNotFound
					}
				}

				return &teacher, err2
			},
		})
	}, nil)

	if err == ErrorObjectNotFound || err == nil && teacher.ID == 0 {
		return nil, ErrorObjectNotFound
	}

	return &teacher, err
}
