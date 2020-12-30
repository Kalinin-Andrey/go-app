package repository

import (
	"context"

	"github.com/stretchr/testify/mock"

	"redditclone/internal/domain"
	"redditclone/internal/domain/post"
)

// UserRepository is a mock for UserRepository
type PostRepository struct {
	mock.Mock
}

var _ post.Repository = (*PostRepository)(nil)

func (m PostRepository) SetDefaultConditions(conditions domain.DBQueryConditions) {}

func (m PostRepository) Get(a0 context.Context, a1 string) (*post.Post, error) {
	ret := m.Called(a0, a1)

	var r0 *post.Post
	if rf, ok := ret.Get(0).(func(context.Context, string) *post.Post); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*post.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(a0, a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m PostRepository) Query(a0 context.Context, a1 domain.DBQueryConditions) ([]post.Post, error) {
	ret := m.Called(a0, a1)

	var r0 []post.Post
	if rf, ok := ret.Get(0).(func(context.Context, domain.DBQueryConditions) []post.Post); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]post.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.DBQueryConditions) error); ok {
		r1 = rf(a0, a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m PostRepository) Create(a0 context.Context, a1 *post.Post) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *post.Post) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m PostRepository) Update(a0 context.Context, a1 *post.Post) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *post.Post) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m PostRepository) Delete(a0 context.Context, a1 string) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
