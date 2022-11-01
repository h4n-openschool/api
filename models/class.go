package models

import "time"

type Class struct {
	Id string `json:"id"`

	Name        string `json:"name"`
	DisplayName string `json:"display_name"`

	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
