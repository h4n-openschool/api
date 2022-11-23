package handlers

import (
	"github.com/h4n-openschool/api/bus"
	"github.com/h4n-openschool/api/repos/classes"
	"github.com/h4n-openschool/api/repos/teachers"
	"go.uber.org/zap"
)

// OpenSchoolImpl implements the [api.ServerInterface] to implement the contract
// defined by the OpenAPI specification (and the generated interfaces from it).
type OpenSchoolImpl struct {
	ClassRepository   classes.ClassRepository
	TeacherRepository teachers.TeacherRepository
	Bus               *bus.Bus
	Logger            *zap.Logger
}
