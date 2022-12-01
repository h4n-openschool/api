package models

import (
	"time"

	"github.com/h4n-openschool/api/api"
)

// Class represents a class students can attend
type Class struct {
	// Name is the computer-friendly name of the class
	Name string `json:"name"`

	// DisplayName is the human-friendly name of the class
	DisplayName string `json:"display_name"`

	// Description is the human-readable description of what the class teaches.
	Description *string `json:"description"`

	StudentIds []string `json:"studentIds"`

	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	BaseMetadata
}

func (c *Class) AsApiClass() api.Class {
	return api.Class{
		Id:          c.Id,
		Name:        c.Name,
		DisplayName: c.DisplayName,
		Description: *c.Description,
		StudentIds:  &c.StudentIds,
		StartDate:   c.StartDate.Format(time.RFC3339),
		EndDate:     c.EndDate.Format(time.RFC3339),
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
