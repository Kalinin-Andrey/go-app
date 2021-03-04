package post

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"redditclone/internal/domain/comment"
	"redditclone/internal/domain/user"
	"redditclone/internal/domain/vote"
)

const (
	EntityName = "post"
	TableName  = "post"
	TypeText   = "text"
	TypeLink   = "link"

	CategoryMusic       = "music"
	CategoryFunny       = "funny"
	CategoryVideos      = "videos"
	CategoryProgramming = "programming"
	CategoryNews        = "news"
	CategoryFashion     = "fashion"
)

var Types []interface{} = []interface{}{
	TypeText,
	TypeLink,
}

var Categories []string = []string{
	CategoryMusic,
	CategoryFunny,
	CategoryVideos,
	CategoryProgramming,
	CategoryNews,
	CategoryFashion,
}

// Post is the user entity
type Post struct {
	ID       string `gorm:"PRIMARY_KEY" json:"id"`
	Score    int    `json:"score"`
	Views    uint   `json:"views"`
	Title    string `gorm:"type:varchar(100)" json:"title"`
	Type     string `gorm:"type:varchar(100)" json:"type"`
	Category string `gorm:"type:varchar(100)" json:"category"`
	Text     string `json:"text,omitempty"`
	Link     string `gorm:"type:varchar(100)" json:"link,omitempty"`

	UserID uint      `sql:"type:int REFERENCES \"user\"(id)" json:"userId"`
	User   user.User `gorm:"FOREIGNKEY:UserID;association_autoupdate:false" json:"author"`

	Votes    []vote.Vote       `gorm:"FOREIGNKEY:PostID" json:"votes"`
	Comments []comment.Comment `gorm:"FOREIGNKEY:PostID" json:"comments"`

	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	DeletedAt *time.Time `gorm:"INDEX" json:"deleted"`
}

func (e Post) Validate() error {

	err := validation.ValidateStruct(&e,
		validation.Field(&e.Type, validation.Required, validation.Length(2, 100), is.Alpha, validation.In(Types...)),
		validation.Field(&e.Category, validation.Required, validation.Length(2, 100), is.Alpha),
		validation.Field(&e.Title, validation.Required, validation.Length(2, 100)),
	)
	if err != nil {
		return err
	}

	switch e.Type {
	case TypeText:
		err = e.validateText()
	case TypeLink:
		err = e.validateLink()
	}
	return err
}

func (e Post) validateText() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Text, validation.Required),
	)
}

func (e Post) validateLink() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Link, validation.Required, validation.Length(2, 100), is.URL),
	)
}

func (e Post) TableName() string {
	return TableName
}

// New func is a constructor for the Post
func New() *Post {
	return &Post{}
}
