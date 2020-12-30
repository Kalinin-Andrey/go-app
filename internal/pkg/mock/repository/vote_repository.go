package repository

import (
	"context"

	"github.com/stretchr/testify/mock"

	"redditclone/internal/domain"
	"redditclone/internal/domain/vote"
)

// UserRepository is a mock for UserRepository
type VoteRepository struct {
	mock.Mock
}

var _ vote.Repository = (*VoteRepository)(nil)

func (m VoteRepository) SetDefaultConditions(conditions domain.DBQueryConditions) {}

func (m VoteRepository) Get(a0 context.Context, a1 string) (*vote.Vote, error) {
	ret := m.Called(a0, a1)

	var r0 *vote.Vote
	if rf, ok := ret.Get(0).(func(context.Context, string) *vote.Vote); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*vote.Vote)
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

func (m VoteRepository) Query(a0 context.Context, a1 domain.DBQueryConditions) ([]vote.Vote, error) {
	ret := m.Called(a0, a1)

	var r0 []vote.Vote
	if rf, ok := ret.Get(0).(func(context.Context, domain.DBQueryConditions) []vote.Vote); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]vote.Vote)
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

func (m VoteRepository) First(a0 context.Context, a1 *vote.Vote) (*vote.Vote, error) {
	ret := m.Called(a0, a1)

	var r0 *vote.Vote
	if rf, ok := ret.Get(0).(func(context.Context, *vote.Vote) *vote.Vote); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*vote.Vote)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *vote.Vote) error); ok {
		r1 = rf(a0, a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m VoteRepository) Create(a0 context.Context, a1 *vote.Vote) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *vote.Vote) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m VoteRepository) Update(a0 context.Context, a1 *vote.Vote) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *vote.Vote) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m VoteRepository) Delete(a0 context.Context, a1 string) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
