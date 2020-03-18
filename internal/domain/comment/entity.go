package comment

import (
	"time"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
)

const (
	TableName	= "comment"
)

// Entity is the user entity
type Entity struct {
	ID				uint			`gorm:"PRIMARY_KEY" json:"id"`
	PostID			uint			`gorm:"" json:"postId"`
	UserID			uint			`gorm:"" json:"userId"`
	User			*user.Entity	`gorm:"EMBEDDED" json:"author"`
	Body			string			`json:"body"`

	CreatedAt		time.Time		`json:"created"`
	UpdatedAt		time.Time		`json:"updated"`
	DeletedAt		*time.Time	`gorm:"INDEX" json:"deleted"`
}


func (e Entity) TableName() string {
	return TableName
}

// New func is a constructor for the Entity
func New() *Entity {
	return &Entity{}
}

