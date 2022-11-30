package handlers

import (
	"github.com/h4n-openschool/api/repos/classes"
	"github.com/h4n-openschool/api/repos/students"
	"github.com/h4n-openschool/api/repos/teachers"
	"go.uber.org/zap"
)

// OpenSchoolImpl implements the [api.ServerInterface] to implement the contract
// defined by the OpenAPI specification (and the generated interfaces from it).
type OpenSchoolImpl struct {
	StudentRepository students.StudentRepository
	ClassRepository   classes.ClassRepository
	TeacherRepository teachers.TeacherRepository
	Logger            *zap.Logger
}
