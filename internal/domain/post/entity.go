package post

import (
	"github.com/Kalinin-Andrey/redditclone/internal/domain/comment"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
	"github.com/Kalinin-Andrey/redditclone/internal/domain/vote"
	"time"
)

const (
	TableName	= "post"
	TypeText	= "text"
	TypeLink	= "link"

	CategoryMusic		= "music"
	CategoryFunny		= "funny"
	CategoryVideos		= "videos"
	CategoryProgramming	= "programming"
	CategoryNews		= "news"
	CategoryFashion		= "fashion"
)

var Types []string = []string{
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
	ID				uint		`gorm:"PRIMARY_KEY" json:"id"`
	Score			uint		`json:"score"`
	Views			uint		`json:"views"`
	Title			string		`gorm:"type:varchar(100)" json:"title"`
	Type			string		`gorm:"type:varchar(100)" json:"type"`
	Category		string		`gorm:"type:varchar(100)" json:"category"`
	Text			string		`json:"text,omitempty"`
	Link			string		`gorm:"type:varchar(100)" json:"link,omitempty"`

	UserID			uint     	`sql:"type:int REFERENCES \"user\"(id)" json:"userId"`
	User			user.User	`gorm:"FOREIGNKEY:UserID" json:"author"`

	Votes			[]vote.Vote      `gorm:"FOREIGNKEY:PostID" json:"votes"`
	Comments		[]comment.Comment `gorm:"FOREIGNKEY:PostID" json:"comments"`

	CreatedAt		time.Time		`json:"created"`
	UpdatedAt		time.Time		`json:"updated"`
	DeletedAt		*time.Time	`gorm:"INDEX" json:"deleted"`
}


func (e Post) TableName() string {
	return TableName
}

// New func is a constructor for the Post
func New() *Post {
	return &Post{}
}

