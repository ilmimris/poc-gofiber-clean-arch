package psql

import (
	"context"
	"database/sql"

	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain"
)

type psqlAuthorRepo struct {
	DB *sql.DB
}

// NewPsqlAuthorRepository will create an implementation of author repository
func NewPsqlAuthorRepository(db *sql.DB) domain.AuthorRepository {
	return &psqlAuthorRepo{
		DB: db,
	}
}

func (p *psqlAuthorRepo) getOne(ctx context.Context, query string, args ...interface{}) (res domain.Author, err error) {
	statement, err := p.DB.PrepareContext(ctx, query)

	// if error
	if err != nil {
		return domain.Author{}, err
	}

	row := statement.QueryRowContext(ctx, args...)
	res = domain.Author{}

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	return
}

func (p *psqlAuthorRepo) GetByID(ctx context.Context, id int64) (domain.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM public.author WHERE id=$1`
	return p.getOne(ctx, query, id)
}
