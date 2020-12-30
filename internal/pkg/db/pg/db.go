package pg

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"redditclone/internal/pkg/config"
	"redditclone/internal/pkg/log"
)

// IDB is the interface for a DB connection
type IDB interface {
	DB() *gorm.DB
}

// DB is the struct for a DB connection
type DB struct {
	D *gorm.DB
}

func (db *DB) DB() *gorm.DB {
	return db.D
}

var _ IDB = (*DB)(nil)

// New creates a new DB connection
func New(conf config.Pg, logger log.ILogger) (*DB, error) {
	db, err := gorm.Open(conf.Dialect, conf.DSN)
	if err != nil {
		return nil, err
	}
	db.SetLogger(logger)
	// Enable Logger, show detailed log
	db.LogMode(conf.IsLogMode)
	// Enable auto preload embeded entities
	db = db.Set("gorm:auto_preload", true)

	dbobj := &DB{D: db}

	if conf.IsAutoMigrate {
		dbobj.AutoMigrateAll()
	}

	return dbobj, nil
}
