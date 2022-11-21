package handlers

import (
	"github.com/h4n-openschool/api/bus"
	"github.com/h4n-openschool/api/repos/classes"
)

// OpenSchoolImpl implements the [api.ServerInterface] to implement the contract
// defined by the OpenAPI specification (and the generated interfaces from it).
type OpenSchoolImpl struct {
	Repository classes.ClassRepository
	Bus        *bus.Bus
}
