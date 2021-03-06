package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/common/repository"
	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain"
	postRepo "github.com/ilmimris/poc-gofiber-clean-arch/pkg/post/repository/mysql"
	"github.com/stretchr/testify/assert"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockPost := []domain.Post{
		domain.Post{
			ID: 1, Title: "title 1", Content: "content 1",
			Author: domain.Author{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		domain.Post{
			ID: 2, Title: "title 2", Content: "content 2",
			Author: domain.Author{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(mockPost[0].ID, mockPost[0].Title, mockPost[0].Content,
			mockPost[0].Author.ID, mockPost[0].UpdatedAt, mockPost[0].CreatedAt).
		AddRow(mockPost[1].ID, mockPost[1].Title, mockPost[1].Content,
			mockPost[1].Author.ID, mockPost[1].UpdatedAt, mockPost[1].CreatedAt)

	query := "SELECT id, title, content, author_id, updated_at, created_at FROM post WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	entry := postRepo.NewMysqlPostRepository(db)
	cursor := repository.EncodeCursor(mockPost[1].CreatedAt)
	num := int64(2)

	list, nextCursor, err := entry.Fetch(context.TODO(), cursor, num)

	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)

}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id, title, content, author_id, updated_at, created_at FROM post WHERE id = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	entry := postRepo.NewMysqlPostRepository(db)

	num := int64(5)
	anPost, err := entry.GetByID(context.TODO(), num)

	assert.NoError(t, err)
	assert.NotNil(t, anPost)
}

func TestStore(t *testing.T) {
	now := time.Now()
	post := &domain.Post{
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: domain.Author{
			ID:   1,
			Name: "Dummy User",
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT post SET title=\\? , content=\\? , author_id=\\? , updated_at=\\? , created_at=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(post.Title, post.Content, post.Author.ID, post.UpdatedAt, post.CreatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

	entry := postRepo.NewMysqlPostRepository(db)
	err = entry.Store(context.TODO(), post)

	assert.NoError(t, err)
	assert.Equal(t, int64(12), post.ID)
}

func TestGetByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "updated_at", "created_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now())

	query := "SELECT id, title, content, author_id, updated_at, created_at FROM post WHERE title = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	entry := postRepo.NewMysqlPostRepository(db)

	title := "title 1"
	anPost, err := entry.GetByTitle(context.TODO(), title)

	assert.NoError(t, err)
	assert.NotNil(t, anPost)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM post WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	entry := postRepo.NewMysqlPostRepository(db)

	num := int64(12)
	err = entry.Delete(context.TODO(), num)

	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	post := &domain.Post{
		ID:        12,
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: domain.Author{
			ID:   1,
			Name: "Dummy User",
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE post set title=\\?, content=\\?, author_id=\\?, updated_at=\\? WHERE ID = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(post.Title, post.Content, post.Author.ID, post.UpdatedAt, post.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	entry := postRepo.NewMysqlPostRepository(db)

	err = entry.Update(context.TODO(), post)

	assert.NoError(t, err)
}
