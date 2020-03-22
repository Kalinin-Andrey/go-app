package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/Kalinin-Andrey/redditclone/pkg/config"
	"github.com/Kalinin-Andrey/redditclone/pkg/log"

	//"database/sql"
	//_ "github.com/lib/pq"
)

// IDB is the interface for a DB connection
type IDB interface {
	DB() *gorm.DB
}

// DB is the struct for a DB connection
type DB struct {
	db *gorm.DB
}

func (db *DB) DB() *gorm.DB {
	return db.db
}

var _ IDB = (*DB)(nil)

// New creates a new DB connection
func New(conf config.DB, logger log.ILogger) (*DB, error) {
	db, err := gorm.Open(conf.Dialect, conf.DSN)
	if err != nil {
		return nil, err
	}
	db.SetLogger(logger)
	// Enable Logger, show detailed log
	db.LogMode(true)
	// Enable auto preload embeded entities
	db = 	db.Set("gorm:auto_preload", true)

	dbobj := &DB{db: db}
	dbobj.AutoMigrateAll()
	
	return dbobj, nil
}

