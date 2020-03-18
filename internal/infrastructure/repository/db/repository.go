package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/redditclone/pkg/db"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"
)

// IRepository is an interface of repository
type IRepository interface {}

// repository persists albums in database
type repository struct {
	db     db.IDB
	logger log.ILogger
	defaultConditions	map[string]interface{}
}

const Limit = 100

// GetRepository return a repository
func GetRepository(dbase db.IDB, logger log.ILogger, entity string) (repo IRepository, err error) {
	r := &repository{
		db:     dbase,
		logger: logger,
	}

	switch entity {
	case "user":
		repo, err = NewUserRepository(r)
	case "post":
		repo, err = NewPostRepository(r)
	case "vote":
		repo, err = NewVoteRepository(r)
	case "comment":
		repo, err = NewCommentRepository(r)
	default:
		err = errors.Errorf("Repository for entity %q not found", entity)
	}
	return repo, err
}


func  (r *repository) SetDefaultConditions(conditions map[string]interface{}) {
	r.defaultConditions = conditions

	if _, ok := r.defaultConditions["Limit"]; !ok {
		r.defaultConditions["Limit"] = Limit
	}
}

func (r repository) dbWithDefaults() *gorm.DB {
	db := r.db.DB()

	if where, ok := r.defaultConditions["Where"]; ok {
		db = db.Where(where)
	}

	if order, ok := r.defaultConditions["SortOrder"]; ok {
		db = db.Order(order)
	}

	if limit, ok := r.defaultConditions["Limit"]; ok {
		db = db.Limit(limit)
	}

	return db
}

