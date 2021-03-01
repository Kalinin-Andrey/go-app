package redis

import (
	"github.com/minipkg/go-app-common/db/redis"
	"redditclone/internal/domain"
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
