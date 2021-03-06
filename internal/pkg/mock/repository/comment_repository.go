package repository

import (
	"context"

	"github.com/stretchr/testify/mock"

	"redditclone/internal/domain"
	"redditclone/internal/domain/comment"
)

// UserRepository is a mock for UserRepository
type CommentRepository struct {
	mock.Mock
}

var _ comment.Repository = (*CommentRepository)(nil)

func (m CommentRepository) SetDefaultConditions(conditions domain.DBQueryConditions) {}

func (m CommentRepository) Get(a0 context.Context, a1 string) (*comment.Comment, error) {
	ret := m.Called(a0, a1)

	var r0 *comment.Comment
	if rf, ok := ret.Get(0).(func(context.Context, string) *comment.Comment); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*comment.Comment)
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

func (m CommentRepository) Query(a0 context.Context, a1 domain.DBQueryConditions) ([]comment.Comment, error) {
	ret := m.Called(a0, a1)

	var r0 []comment.Comment
	if rf, ok := ret.Get(0).(func(context.Context, domain.DBQueryConditions) []comment.Comment); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]comment.Comment)
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

func (m CommentRepository) Create(a0 context.Context, a1 *comment.Comment) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *comment.Comment) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m CommentRepository) Update(a0 context.Context, a1 *comment.Comment) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *comment.Comment) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m CommentRepository) Delete(a0 context.Context, a1 string) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
