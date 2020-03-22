package db

import (
	"github.com/Kalinin-Andrey/redditclone/internal/domain/comment"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/post"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/session"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/vote"
)

func (db *DB) AutoMigrateAll() {
	db.DB().AutoMigrate(
		&user.User{},
		&session.Session{},
		&post.Post{},
		&vote.Vote{},
		&comment.Comment{},
		)
}
