package domain

import (
	"context"
	"time"
)

// Post represent the post model
type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Author    Author    `json:"author"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// PostUsecase represent the post's usecase contract
type PostUsecase interface {
	// Create
	Store(ctx context.Context, p *Post) error

	// Read
	Fetch(ctx context.Context, cursor string, num int64) ([]Post, string, error)
	GetByID(ctx context.Context, id int64) (Post, error)
	GetByTitle(ctx context.Context, title string) (Post, error)

	// Update
	Update(ctx context.Context, p *Post) error

	// Delete
	Delete(ctx context.Context, id int64) error
}

// PostRepository represent the post's repository contract
type PostRepository interface {
	// Create
	Store(ctx context.Context, p *Post) error

	// Read
	Fetch(ctx context.Context, cursor string, num int64) (res []Post, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (Post, error)
	GetByTitle(ctx context.Context, title string) (Post, error)

	// Update
	Update(ctx context.Context, p *Post) error

	// Delete
	Delete(ctx context.Context, id int64) (err error)
}
