package models

import "time"

// Post represent the post model
type Post struct {
	ID        uint32    `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Author    Author    `json:"author"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
