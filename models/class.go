package models

import (
	"time"

	"github.com/h4n-openschool/api/api"
)

// Class represents a class students can attend
type Class struct {
	BaseMetadata

	// Name is the computer-friendly name of the class
	Name string `json:"name"`

	// DisplayName is the human-friendly name of the class
	DisplayName string `json:"display_name"`

	// Description is the human-readable description of what the class teaches.
	Description *string `json:"description"`
}

func (c *Class) AsApiClass() api.Class {
	return api.Class{
		Id:          c.Id,
		Name:        c.Name,
		DisplayName: c.DisplayName,
		Description: *c.Description,
		CreatedAt:   c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   c.UpdatedAt.Format(time.RFC3339),
	}
}

func (c *Class) ReconcileWithApiClass(description *string, displayName *string) *Class {
	if description != nil {
		c.Description = description
	}

	if displayName != nil {
		c.DisplayName = *displayName
	}

	return c
}

func ClassesAsApiClassList(classes []Class) api.ClassList {
	classList := api.ClassList{}
	for _, class := range classes {
		classList = append(classList, class.AsApiClass())
	}
	return classList
}
