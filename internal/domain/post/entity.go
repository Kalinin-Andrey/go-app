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

// Entity is the user entity
type Entity struct {
	ID				uint		`gorm:"PRIMARY_KEY" json:"id"`
	Score			uint		`json:"score"`
	Views			uint		`json:"views"`
	Title			string		`gorm:"type:varchar(100)" json:"title"`
	Type			string		`gorm:"type:varchar(100)" json:"type"`
	Category		string		`gorm:"type:varchar(100)" json:"category"`
	Text			string		`json:"text,omitempty"`
	Link			string		`gorm:"type:varchar(100)" json:"link,omitempty"`

	UserID			uint			`json:"userId"`
	User			*user.Entity	`gorm:"EMBEDDED" json:"author"`

	Votes			[]vote.Entity	`gorm:"EMBEDDED" json:"votes"`
	Comments		[]comment.Entity	`gorm:"EMBEDDED" json:"comments"`

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

