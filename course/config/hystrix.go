package config

import (
	"github.com/afex/hystrix-go/hystrix"
)

func init() {
	hystrix.Configure(map[string]hystrix.CommandConfig{
		"CourseCUD": {
			Timeout:                5000,
			MaxConcurrentRequests:  20,
			RequestVolumeThreshold: 20,
			SleepWindow:            5000,
			ErrorPercentThreshold:  25,
		},
		"TeacherCUD": {
			Timeout:                5000,
			MaxConcurrentRequests:  20,
			RequestVolumeThreshold: 20,
			SleepWindow:            5000,
			ErrorPercentThreshold:  25,
		},
		"GetAllCourses": {
			Timeout:                2000,
			MaxConcurrentRequests:  20,
			RequestVolumeThreshold: 20,
			SleepWindow:            5000,
			ErrorPercentThreshold:  25,
		},
		"GetCourseById": {
			Timeout:                1000,
			MaxConcurrentRequests:  50,
			RequestVolumeThreshold: 20,
			SleepWindow:            5000,
			ErrorPercentThreshold:  25,
		},
		"GetAllTeachers": {
			Timeout:                2000,
			MaxConcurrentRequests:  20,
			RequestVolumeThreshold: 20,
			SleepWindow:            5000,
			ErrorPercentThreshold:  25,
		},
		"GetTeacherById": {
			Timeout:                1000,
			MaxConcurrentRequests:  50,
			RequestVolumeThreshold: 20,
			SleepWindow:            5000,
			ErrorPercentThreshold:  25,
		},
		"GetCoursesByTeacherId": {
			Timeout:                1000,
			MaxConcurrentRequests:  50,
			RequestVolumeThreshold: 20,
			SleepWindow:            5000,
			ErrorPercentThreshold:  25,
		},
		"GetTeacherByCourseId": {
			Timeout:                3000,
			MaxConcurrentRequests:  50,
			RequestVolumeThreshold: 20,
			SleepWindow:            5000,
			ErrorPercentThreshold:  25,
		},
	})
}
