package domain

import (
	"context"
	"time"
)

// Author repesent the author Model
type Author struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthorUsecase represent the author's usecase contract
type AuthorUsecase interface {
	// Read
	GetByID(ctx context.Context, id int64) (Author, error)
}

// AuthorRepository represent the author's repository contract
type AuthorRepository interface {
	// Read
	GetByID(ctx context.Context, id int64) (Author, error)
}
