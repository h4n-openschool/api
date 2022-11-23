package models

import "time"

type BaseMetadata struct {
	// Id is the cuid of the class
	Id string `json:"id"`

  // CreatedAt is the time at which this class was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the time at which this class was updated.
	UpdatedAt time.Time `json:"updated_at"`
}

