package models

import (
	"time"

	"github.com/h4n-openschool/api/api"
)

type Grade struct {
  BaseMetadata

  ClassId string `json:"classId"`

  // StudentId is the ID of the student this Grade applies to.
  StudentId string `json:"studentId"`

  // Value is the value of the Grade.
  Value int `json:"value"`
}

func (s *Grade) AsApiGrade() api.Grade {
  return api.Grade{
    Id: s.BaseMetadata.Id,
    CreatedAt: s.BaseMetadata.CreatedAt.Format(time.RFC3339),
    UpdatedAt: s.BaseMetadata.UpdatedAt.Format(time.RFC3339),
    StudentId: s.StudentId,
    Value: s.Value,
  }
}

func GradesAsApiGradeList(students []Grade) api.GradeList {
	studentList := api.GradeList{}
	for _, student := range students {
		studentList = append(studentList, student.AsApiGrade())
	}
	return studentList
}

type Student struct {
  BaseMetadata

  // FullName is the student's full legal name.
  FullName string `json:"fullName"`

  ClassId string `json:"classId"`
}

func (s *Student) AsApiStudent() api.Student {
  return api.Student{
    Id: s.BaseMetadata.Id,
    CreatedAt: s.BaseMetadata.CreatedAt.Format(time.RFC3339),
    UpdatedAt: s.BaseMetadata.UpdatedAt.Format(time.RFC3339),
    FullName: s.FullName,
    ClassId: &s.ClassId,
  }
}

func StudentsAsApiStudentList(students []Student) api.StudentList {
	studentList := api.StudentList{}
	for _, student := range students {
		studentList = append(studentList, student.AsApiStudent())
	}
	return studentList
}
