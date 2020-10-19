package psql_test

import (
	"context"
	"testing"
	"time"

	authorRepo "github.com/ilmimris/poc-gofiber-clean-arch/pkg/author/repository/psql"
	"github.com/stretchr/testify/assert"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(1, "Dummy User", time.Now(), time.Now())

	query := "SELECT id, name, created_at, updated_at FROM public.author WHERE id=\\$1"

	userID := int64(1)
	prep := mock.ExpectPrepare(query)
	prep.ExpectQuery().WithArgs(userID).WillReturnRows(rows)

	a := authorRepo.NewPsqlAuthorRepository(db)

	anArticle, err := a.GetByID(context.TODO(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, anArticle)

}
