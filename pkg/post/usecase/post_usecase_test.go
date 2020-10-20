package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain"
	"github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain/mocks"
	ucase "github.com/ilmimris/poc-gofiber-clean-arch/pkg/post/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetch(t *testing.T) {
	mockPostRepo := new(mocks.PostRepository)
	mockPost := domain.Post{
		Title:   "Hello",
		Content: "Content",
	}

	mockListArtilce := make([]domain.Post, 0)
	mockListArtilce = append(mockListArtilce, mockPost)

	t.Run("success", func(t *testing.T) {
		mockPostRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListArtilce, "next-cursor", nil).Once()
		mockAuthor := domain.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(mocks.AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := ucase.NewPostUsecase(mockPostRepo, mockAuthorrepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArtilce))

		mockPostRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockPostRepo.On("Fetch", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewPostUsecase(mockPostRepo, mockAuthorrepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.Fetch(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockPostRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestGetByID(t *testing.T) {
	mockPostRepo := new(mocks.PostRepository)
	mockPost := domain.Post{
		Title:   "Hello",
		Content: "Content",
	}
	mockAuthor := domain.Author{
		ID:   1,
		Name: "Iman Tumorang",
	}

	t.Run("success", func(t *testing.T) {
		mockPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockPost, nil).Once()
		mockAuthorrepo := new(mocks.AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := ucase.NewPostUsecase(mockPostRepo, mockAuthorrepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockPost.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockPostRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Post{}, errors.New("Unexpected")).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewPostUsecase(mockPostRepo, mockAuthorrepo, time.Second*2)

		a, err := u.GetByID(context.TODO(), mockPost.ID)

		assert.Error(t, err)
		assert.Equal(t, domain.Post{}, a)

		mockPostRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestStore(t *testing.T) {
	mockPostRepo := new(mocks.PostRepository)
	mockPost := domain.Post{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		tempMockPost := mockPost
		tempMockPost.ID = 0
		mockPostRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(domain.Post{}, domain.ErrNotFound).Once()
		mockPostRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.Post")).Return(nil).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewPostUsecase(mockPostRepo, mockAuthorrepo, time.Second*2)

		err := u.Store(context.TODO(), &tempMockPost)

		assert.NoError(t, err)
		assert.Equal(t, mockPost.Title, tempMockPost.Title)
		mockPostRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingPost := mockPost
		mockPostRepo.On("GetByTitle", mock.Anything, mock.AnythingOfType("string")).Return(existingPost, nil).Once()
		mockAuthor := domain.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(mocks.AuthorRepository)
		mockAuthorrepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)

		u := ucase.NewPostUsecase(mockPostRepo, mockAuthorrepo, time.Second*2)

		err := u.Store(context.TODO(), &mockPost)

		assert.Error(t, err)
		mockPostRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestDelete(t *testing.T) {
	mockPostRepo := new(mocks.PostRepository)
	mockPost := domain.Post{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockPost, nil).Once()

		mockPostRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewPostUsecase(mockPostRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockPost.ID)

		assert.NoError(t, err)
		mockPostRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("post-is-not-exist", func(t *testing.T) {
		mockPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Post{}, nil).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewPostUsecase(mockPostRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockPost.ID)

		assert.Error(t, err)
		mockPostRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockPostRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Post{}, errors.New("Unexpected Error")).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewPostUsecase(mockPostRepo, mockAuthorrepo, time.Second*2)

		err := u.Delete(context.TODO(), mockPost.ID)

		assert.Error(t, err)
		mockPostRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	mockPostRepo := new(mocks.PostRepository)
	mockPost := domain.Post{
		Title:   "Hello",
		Content: "Content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockPostRepo.On("Update", mock.Anything, &mockPost).Once().Return(nil)

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewPostUsecase(mockPostRepo, mockAuthorrepo, time.Second*2)

		err := u.Update(context.TODO(), &mockPost)
		assert.NoError(t, err)
		mockPostRepo.AssertExpectations(t)
	})
}
