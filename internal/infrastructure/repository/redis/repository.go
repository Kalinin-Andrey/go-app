package redis

import (
	"redditclone/internal/domain"

	"github.com/minipkg/db/redis"
)

// IRepository is an interface of repository
type IRepository interface{}

// repository persists albums in database
type repository struct {
	db         redis.IDB
	Conditions domain.DBQueryConditions
}

const DefaultLimit = 100

func (r *repository) SetDefaultConditions(defaultConditions domain.DBQueryConditions) {
	r.Conditions = defaultConditions

	if r.Conditions.Limit == 0 {
		r.Conditions.Limit = DefaultLimit
	}
}
