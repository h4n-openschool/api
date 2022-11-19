package models

import (
	"time"

	"github.com/h4n-openschool/classes/api"
)

// Class represents a class students can attend
type Class struct {
	// Id is the cuid of the class
	Id string `json:"id"`

	// Name is the computer-friendly name of the class
	Name string `json:"name"`

	// DisplayName is the human-friendly name of the class
	DisplayName string `json:"display_name"`

	// Description is the human-readable description of what the class teaches.
	Description string `json:"description"`

	// CreatedAt is the time at which this class was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the time at which this class was updated.
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Class) AsApiClass() api.Class {
  return api.Class{
    Id: c.Id,
    Name: c.Name,
    DisplayName: c.DisplayName,
    Description: c.Description,
    CreatedAt: c.CreatedAt.Format(time.RFC3339),
    UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
  }
}

func ClassesAsApiClassList(classes []Class) (api.ClassList) {
  classList := api.ClassList{}
  for _, class := range classes { 
    classList = append(classList, class.AsApiClass())
  }
  return classList
}

