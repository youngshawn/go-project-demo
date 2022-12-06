package models

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/afex/hystrix-go/hystrix"
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

	IDcard string `json:"-" gorm:"unique"`
	Phone  string `json:"-"`

	PlainIDcard string `json:"idcard" gorm:"-" msgpack:"-"`
	PlainPhone  string `json:"phone" gorm:"-" msgpack:"-"`
}

func (t *Teacher) Encrypt() error {
	if t.IDcard != "" && t.Phone != "" {
		log.Printf("teacher.Encrypt already called, return")
		return nil
	}

	path := fmt.Sprintf("/transit/encrypt/%s", config.VaultTransitKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := Vault.Logical().WriteWithContext(ctx, path, map[string]interface{}{
		"batch_input": []map[string]interface{}{
			{
				"plaintext": base64.StdEncoding.EncodeToString([]byte(t.PlainIDcard)),
			},
			{
				"plaintext": base64.StdEncoding.EncodeToString([]byte(t.PlainPhone)),
			},
		},
	})
	if err != nil {
		log.Printf("teacher.Encrypt failed to encrypt data")
		return err
	}

	t.IDcard = resp.Data["batch_results"].([]interface{})[0].(map[string]interface{})["ciphertext"].(string)
	t.Phone = resp.Data["batch_results"].([]interface{})[1].(map[string]interface{})["ciphertext"].(string)

	log.Printf("teacher.Encrypt called")
	return nil
}

func (t *Teacher) Decrypt() error {
	if t.PlainIDcard != "" && t.PlainPhone != "" {
		log.Printf("teacher.Decrypt already called, return")
		return nil
	}

	path := fmt.Sprintf("/transit/decrypt/%s", config.VaultTransitKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := Vault.Logical().WriteWithContext(ctx, path, map[string]interface{}{
		"batch_input": []map[string]interface{}{
			{
				"ciphertext": t.IDcard,
			},
			{
				"ciphertext": t.Phone,
			},
		},
	})
	if err != nil {
		log.Printf("teacher.Decrypt failed to decrypt data")
		return err
	}

	IDcard_base64 := resp.Data["batch_results"].([]interface{})[0].(map[string]interface{})["plaintext"].(string)
	Phone_base64 := resp.Data["batch_results"].([]interface{})[1].(map[string]interface{})["plaintext"].(string)
	IDcard, err1 := base64.StdEncoding.DecodeString(IDcard_base64)
	Phone, err2 := base64.StdEncoding.DecodeString(Phone_base64)
	if err1 != nil || err2 != nil {
		log.Printf("teacher.Decrypt failed to base64 decode")
		return errors.New("base64 decode error")
	}
	t.PlainIDcard = string(IDcard)
	t.PlainPhone = string(Phone)

	log.Printf("teacher.Decrypt called")
	return nil
}

func (t *Teacher) BeforeSave(*gorm.DB) error {
	return t.Encrypt()
}

func (t *Teacher) CreateTeacher() error {

	err := hystrix.Do("TeacherCUD", func() error {

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

	}, nil)

	return err
}

func (t *Teacher) UpdateTeacher() (*Teacher, error) {
	var teacher Teacher

	err := hystrix.Do("TeacherCUD", func() error {

		result := db.Find(&teacher, t.ID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrorObjectNotFound
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

		result = db.Omit("CreatedAt").Save(t)

		time.Sleep(time.Millisecond * 500)
		Cache.DeleteFromLocalCache("/teacher/")
		Cache.Delete(context.Background(), "/teacher/")
		Cache.DeleteFromLocalCache(fmt.Sprintf("/teacher/%d", t.ID))
		Cache.Delete(context.Background(), fmt.Sprintf("/teacher/%d", t.ID))
		for _, c := range courses {
			Cache.DeleteFromLocalCache(fmt.Sprintf("/course/%d/teacher", c.ID))
			Cache.Delete(context.Background(), fmt.Sprintf("/course/%d/teacher", c.ID))
		}

		return result.Error

	}, nil)

	if err != nil {
		return nil, err
	}
	return t, nil
}

func DeleteTeacherById(Id uint) (*Teacher, error) {
	var teacher Teacher

	err := hystrix.Do("TeacherCUD", func() error {

		result := db.Find(&teacher, Id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrorObjectNotFound
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

		return result.Error

	}, nil)

	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

func GetAllTeachers() ([]Teacher, error) {
	dynamicCacheConfig := config.GetDynamicCacheConfig()

	var teachers []Teacher

	err := hystrix.Do("GetAllTeachers", func() error {
		return Cache.Once(&cache.Item{
			Key:            "/teacher/",
			Value:          &teachers,
			TTL:            dynamicCacheConfig.CacheTTL,
			SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
			Do: func(i *cache.Item) (interface{}, error) {
				log.Println("GetAllTeachers: miss cache get from DB...")
				result := db.Find(i.Value)

				err := result.Error
				if err == nil && result.RowsAffected == 0 && !dynamicCacheConfig.EnableNullResultCache {
					err = ErrorObjectNotFound
				}

				return i.Value, err
			},
		})
	}, nil)

	if err == nil && len(teachers) == 0 {
		err = ErrorObjectNotFound
	}

	for idx := range teachers {
		teachers[idx].Decrypt()
	}

	return teachers, err
}

func GetTeacherById(Id uint) (*Teacher, error) {
	dynamicCacheConfig := config.GetDynamicCacheConfig()

	var teacher Teacher

	err := hystrix.Do("GetTeacherById", func() error {
		return Cache.Once(&cache.Item{
			Key:            fmt.Sprintf("/teacher/%d", Id),
			Value:          &teacher,
			TTL:            dynamicCacheConfig.CacheTTL,
			SkipLocalCache: !dynamicCacheConfig.EnableLocalCache,
			Do: func(i *cache.Item) (interface{}, error) {
				log.Println("GetTeacherById: miss cache get from DB...")
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

	if err == ErrorObjectNotFound || err == nil && teacher.ID == 0 {
		return nil, ErrorObjectNotFound
	}

	if err == nil {
		teacher.Decrypt()
	}

	return &teacher, err
}
