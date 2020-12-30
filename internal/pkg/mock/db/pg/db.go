package pg

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"redditclone/internal/pkg/config"
	"redditclone/internal/pkg/db/pg"
	"redditclone/internal/pkg/log"
)

// New creates a new DB connection
func New(conf config.Pg, logger log.ILogger) (*pg.DB, *sqlmock.Sqlmock, error) {
	var mock sqlmock.Sqlmock
	var dbm *sql.DB
	var err error

	dbm, mock, err = sqlmock.New() // mock sql.DB
	/*if err := dbm.Ping(); err != nil {
		return nil, nil, err
	}*/

	db, err := gorm.Open(conf.Dialect, dbm)
	if err != nil {
		return nil, nil, err
	}
	db.SetLogger(logger)
	// Enable Logger, show detailed log
	db.LogMode(conf.IsLogMode)
	// Enable auto preload embeded entities
	db = db.Set("gorm:auto_preload", true)

	dbobj := &pg.DB{D: db}

	return dbobj, &mock, nil
}
