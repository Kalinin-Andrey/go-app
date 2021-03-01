package pg

import (
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"redditclone/internal/domain"
	"redditclone/internal/domain/user"

	"github.com/minipkg/go-app-common/log"

	"github.com/minipkg/go-app-common/db/pg"
)

// IRepository is an interface of repository
type IRepository interface{}

// repository persists albums in database
type repository struct {
	db         pg.IDB
	logger     log.ILogger
	Conditions domain.DBQueryConditions
}

const DefaultLimit = 100

// GetRepository return a repository
func GetRepository(logger log.ILogger, dbase pg.IDB, entity string) (repo IRepository, err error) {
	r := &repository{
		db:     dbase,
		logger: logger,
	}

	switch entity {
	case user.EntityName:
		repo, err = NewUserRepository(r)
	default:
		err = errors.Errorf("Repository for entity %q not found", entity)
	}
	return repo, err
}

func (r *repository) SetDefaultConditions(defaultConditions domain.DBQueryConditions) {
	r.Conditions = defaultConditions

	if r.Conditions.Limit == 0 {
		r.Conditions.Limit = DefaultLimit
	}
}

func (r repository) dbWithDefaults() *gorm.DB {
	db := r.db.DB()

	if r.Conditions.Where != nil {
		m := r.keysToSnakeCase(r.Conditions.Where)
		db = db.Where(m)
	}

	if r.Conditions.SortOrder != nil {
		m := r.keysToSnakeCase(r.Conditions.SortOrder)
		db = db.Order(m)
	}

	if r.Conditions.Limit != 0 {
		db = db.Limit(r.Conditions.Limit)
	}

	return db
}

func (r repository) keysToSnakeCase(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))

	for key, val := range in {
		out[strcase.ToSnake(key)] = val
	}
	return out
}
