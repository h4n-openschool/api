package models

import "time"

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
