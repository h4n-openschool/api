package models

import (
	"time"

	"github.com/h4n-openschool/api/api"
)

type Teacher struct {
	FullName     string `json:"fullName"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	BaseMetadata
}

func (c *Teacher) AsApiTeacher() api.Teacher {
	return api.Teacher{
		Id:        c.Id,
		FullName:  c.FullName,
		Email:     c.Email,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}

func TeachersAsApiTeacherList(teachers []Teacher) api.TeacherList {
	teacherList := api.TeacherList{}
	for _, teacher := range teachers {
		teacherList = append(teacherList, teacher.AsApiTeacher())
	}
	return teacherList
}
