// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/ilmimris/poc-gofiber-clean-arch/pkg/domain"
	mock "github.com/stretchr/testify/mock"
)

// PostRepository is an autogenerated mock type for the PostRepository type
type PostRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *PostRepository) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fetch provides a mock function with given fields: ctx, cursor, num
func (_m *PostRepository) Fetch(ctx context.Context, cursor string, num int64) ([]domain.Post, string, error) {
	ret := _m.Called(ctx, cursor, num)

	var r0 []domain.Post
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) []domain.Post); ok {
		r0 = rf(ctx, cursor, num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Post)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string, int64) string); ok {
		r1 = rf(ctx, cursor, num)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, int64) error); ok {
		r2 = rf(ctx, cursor, num)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *PostRepository) GetByID(ctx context.Context, id int64) (domain.Post, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.Post
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.Post); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Post)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTitle provides a mock function with given fields: ctx, title
func (_m *PostRepository) GetByTitle(ctx context.Context, title string) (domain.Post, error) {
	ret := _m.Called(ctx, title)

	var r0 domain.Post
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Post); ok {
		r0 = rf(ctx, title)
	} else {
		r0 = ret.Get(0).(domain.Post)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, p
func (_m *PostRepository) Store(ctx context.Context, p *domain.Post) error {
	ret := _m.Called(ctx, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Post) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, p
func (_m *PostRepository) Update(ctx context.Context, p *domain.Post) error {
	ret := _m.Called(ctx, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Post) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
