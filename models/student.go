package models

import (
	"time"

	"github.com/h4n-openschool/api/api"
)

type Grade struct {
  BaseMetadata

  // StudentId is the ID of the student this Grade applies to.
  StudentId string `json:"student"`

  // Value is the value of the Grade.
  Value int `json:"value"`
}

type Student struct {
  BaseMetadata

  // Id is the CUID of the student.
  Id string `json:"id"`

  // FullName is the student's full legal name.
  FullName string `json:"fullName"`
}

func (s *Student) AsApiStudent() api.Student {
  return api.Student{
    Id: s.BaseMetadata.Id,
    CreatedAt: s.BaseMetadata.CreatedAt.Format(time.RFC3339),
    UpdatedAt: s.BaseMetadata.UpdatedAt.Format(time.RFC3339),
    FullName: s.FullName,
  }
}

func StudentsAsApiStudentList(students []Student) api.StudentList {
	studentList := api.StudentList{}
	for _, student := range students {
		studentList = append(studentList, student.AsApiStudent())
	}
	return studentList
}
