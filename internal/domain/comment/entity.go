package comment

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
)

const (
	TableName = "comment"
)

// Comment is the user entity
type Comment struct {
	ID     uint      `gorm:"PRIMARY_KEY" json:"id"`
	PostID uint      `sql:"type:int REFERENCES post(id)" json:"postId"`
	UserID uint      `sql:"type:int REFERENCES \"user\"(id)" json:"userId"`
	User   user.User `gorm:"FOREIGNKEY:UserID;association_autoupdate:false" json:"author"`
	Body   string    `json:"body"`

	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	DeletedAt *time.Time `gorm:"INDEX" json:"deleted"`
}

func (e Comment) Validate() error {

	return validation.ValidateStruct(&e,
		validation.Field(&e.Body, validation.Required),
	)
}

func (e Comment) TableName() string {
	return TableName
}

// New func is a constructor for the Comment
func New() *Comment {
	return &Comment{}
}
