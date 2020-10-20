package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/common/repository"
	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain"
)

type mysqlPostRepo struct {
	DB *sql.DB
}

// NewMysqlPostRepository will create an implementation of post repository
func NewMysqlPostRepository(db *sql.DB) domain.PostRepository {
	return &mysqlPostRepo{
		DB: db,
	}
}

func (p *mysqlPostRepo) Store(ctx context.Context, entry *domain.Post) (err error) {
	query := `INSERT post 
				SET title=? , content=? , author_id=? , updated_at=? , created_at=?`

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

func (p *mysqlPostRepo) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Post, err error) {
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

func (p *mysqlPostRepo) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Post, nextCursor string, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at 
				FROM post 
				WHERE created_at > ? 
				ORDER BY created_at 
				LIMIT ?`

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

func (p *mysqlPostRepo) GetByID(ctx context.Context, id int64) (res domain.Post, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at
				FROM post 
				WHERE id = ?`

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

func (p *mysqlPostRepo) GetByTitle(ctx context.Context, title string) (res domain.Post, err error) {
	query := `SELECT id, title, content, author_id, updated_at, created_at
				FROM post 
				WHERE title = ?`

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

func (p *mysqlPostRepo) Update(ctx context.Context, entry *domain.Post) (err error) {
	query := `UPDATE post set title=?, content=?, author_id=?, updated_at=? WHERE ID = ?`

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

func (p *mysqlPostRepo) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM post WHERE id = ?`

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
