package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/common/repository"
	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain"
)

type psqlPostRepo struct {
	DB *sql.DB
}

// NewPsqlPostRepository will create an implementation of post repository
func NewPsqlPostRepository(db *sql.DB) domain.PostRepository {
	return &psqlPostRepo{
		DB: db,
	}
}

func (p *psqlPostRepo) Store(ctx context.Context, entry *domain.Post) (err error) {
	query := `INSERT public.post 
				SET title=$1 , content=$2 , author_id=$3 , updated_at=$4 , created_at=$5`

	statement, err := p.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := statement.ExecContext(ctx, entry.Title, entry.Content, entry.Author.ID, entry.UpdatedAt, entry.CreatedAt)
	if err != nil {
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}

	entry.ID = lastID
	return
}

func (p *psqlPostRepo) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Post, err error) {
	rows, err := p.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Print(errRow)
		}
	}()

	result = make([]domain.Post, 0)
	for rows.Next() {
		t := domain.Post{}
		authorID := int64(0)

		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			log.Print(err)
			return nil, err
		}

		t.Author = domain.Author{
			ID: authorID,
		}
		result = append(result, t)
	}

	return result, nil
}

func (p *psqlPostRepo) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Post, nextCursor string, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at 
				FROM public.post 
				WHERE created_at > $1 
				ORDER BY created_at 
				LIMIT $2`

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = p.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

func (p *psqlPostRepo) GetByID(ctx context.Context, id int64) (res domain.Post, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at
				FROM public.post 
				WHERE id = $1`

	list, err := p.fetch(ctx, query, id)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (p *psqlPostRepo) GetByTitle(ctx context.Context, title string) (res domain.Post, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at
				FROM public.post 
				WHERE title = $1`

	list, err := p.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (p *psqlPostRepo) Update(ctx context.Context, entry *domain.Post) (err error) {
	query := `UPDATE public.post set title=$1, content=$2, author_id=$3, updated_at=$4 WHERE ID = $5`

	statement, err := p.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := statement.ExecContext(ctx, entry.Title, entry.Content, entry.Author.ID, entry.UpdatedAt, entry.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (p *psqlPostRepo) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM public.post WHERE id = $1`

	statement, err := p.DB.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := statement.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowAffected != 1 {
		err = fmt.Errorf("Weird behavior. Total Affected %d", rowAffected)
		return
	}

	return
}
