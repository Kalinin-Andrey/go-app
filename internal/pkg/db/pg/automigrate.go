package pg

import (
	"redditclone/internal/domain/comment"
	"redditclone/internal/domain/post"
	"redditclone/internal/domain/user"
	"redditclone/internal/domain/vote"
	"redditclone/internal/pkg/session"
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
