package author

import (
	"context"
)

// Repository represent the author's repository contract
type Repository interface {
	GetByID(ctx context.Context, id uint32) (*models.author, error)
}
