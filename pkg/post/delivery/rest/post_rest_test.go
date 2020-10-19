package rest_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain"
	mocks "github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain/mocks"
	postRest "github.com/ilmimris/poc-gofiber-clean-arch/pkg/post/delivery/rest"

	"github.com/gofiber/fiber/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	var mockPost domain.Post
	err := faker.FakeData(&mockPost)
	assert.NoError(t, err)
	mockUCase := new(mocks.PostUsecase)
	mockListPost := make([]domain.Post, 0)
	mockListPost = append(mockListPost, mockPost)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(mockListPost, "10", nil)

	e := fiber.New()
	req, err := http.NewRequest("GET", "/posts?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	postRest.NewPostHandler(e, mockUCase)
	rec, err := e.Test(req, -1)
	fmt.Println(rec)

	require.NoError(t, err)

	responseCursor := rec.Header.Get("X-Cursor")
	assert.Equal(t, "10", responseCursor)
	assert.Equal(t, http.StatusOK, rec.StatusCode)
	mockUCase.AssertExpectations(t)
}

func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.PostUsecase)
	num := 1
	cursor := "2"
	mockUCase.On("Fetch", mock.Anything, cursor, int64(num)).Return(nil, "", domain.ErrInternalServerError)

	e := fiber.New()
	req, err := http.NewRequest("GET", "/posts?num=1&cursor="+cursor, strings.NewReader(""))
	assert.NoError(t, err)

	postRest.NewPostHandler(e, mockUCase)
	rec, err := e.Test(req, -1)

	require.NoError(t, err)

	responseCursor := rec.Header.Get("X-Cursor")
	assert.Equal(t, "", responseCursor)
	assert.Equal(t, http.StatusInternalServerError, rec.StatusCode)
	mockUCase.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	var mockPost domain.Post
	err := faker.FakeData(&mockPost)
	assert.NoError(t, err)

	mockUCase := new(mocks.PostUsecase)

	num := int(mockPost.ID)

	mockUCase.On("GetByID", mock.Anything, int64(num)).Return(mockPost, nil)

	e := fiber.New()
	req, err := http.NewRequest("GET", "/posts/"+strconv.Itoa(num), nil)
	assert.NoError(t, err)

	postRest.NewPostHandler(e, mockUCase)
	rec, err := e.Test(req, -1)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.StatusCode)
	mockUCase.AssertExpectations(t)
}

func TestStore(t *testing.T) {
	mockPost := domain.Post{
		Title:     "Title",
		Content:   "Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tempMockPost := mockPost
	tempMockPost.ID = 0
	mockUCase := new(mocks.PostUsecase)

	j, err := json.Marshal(tempMockPost)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*domain.Post")).Return(nil)

	e := fiber.New()
	req, err := http.NewRequest("POST", "/posts", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	postRest.NewPostHandler(e, mockUCase)
	rec, err := e.Test(req, -1)

	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.StatusCode)
	mockUCase.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	var mockPost domain.Post
	err := faker.FakeData(&mockPost)
	assert.NoError(t, err)

	mockUCase := new(mocks.PostUsecase)

	num := int(mockPost.ID)

	mockUCase.On("Delete", mock.Anything, int64(num)).Return(nil)

	e := fiber.New()
	req, err := http.NewRequest("DELETE", "/posts/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	postRest.NewPostHandler(e, mockUCase)
	rec, err := e.Test(req, -1)

	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.StatusCode)
	mockUCase.AssertExpectations(t)

}
