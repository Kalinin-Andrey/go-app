package vote

import (
	"time"

	"redditclone/internal/domain/user"
)

const (
	EntityName = "vote"
	TableName  = "vote"
)

// Vote is the user entity
type Vote struct {
	ID     string    `gorm:"PRIMARY_KEY" json:"id"`
	PostID string    `sql:"type:varchar REFERENCES post(id)" json:"postId"`
	UserID uint      `sql:"type:int REFERENCES \"user\"(id)" json:"user"`
	User   user.User `gorm:"FOREIGNKEY:UserID;association_autoupdate:false" json:"author"`
	Value  int       `json:"vote"`

	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	DeletedAt *time.Time `gorm:"INDEX" json:"deleted"`
}

func (e Vote) TableName() string {
	return TableName
}

// New func is a constructor for the Vote
func New() *Vote {
	return &Vote{}
}
