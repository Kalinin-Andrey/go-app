package session

import (
	"context"
	"encoding"
	"time"

	"redditclone/internal/pkg/proto"

	"redditclone/internal/domain/user"
)

const (
	EntityName = "session"
	TableName  = "session"
)

type Data struct {
	UserID              uint
	UserName            string
	ExpirationTokenTime time.Time
}

// Session is the session entity
type Session struct {
	ID     uint            `gorm:"PRIMARY_KEY" json:"id"`
	UserID uint            `sql:"type:int NOT NULL REFERENCES \"user\"(id)" json:"userId"`
	User   user.User       `gorm:"FOREIGNKEY:UserID;association_autoupdate:false" json:"author"`
	Data   Data            `gorm:"-"`
	Ctx    context.Context `gorm:"-"`
	Token  string          `gorm:"type:text;unique_index;not null" json:"token"`

	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	DeletedAt *time.Time `gorm:"INDEX" json:"deleted"`
}

var _ encoding.BinaryMarshaler = (*Session)(nil)
var _ encoding.BinaryUnmarshaler = (*Session)(nil)

func (e Session) TableName() string {
	return TableName
}

// New func is a constructor for the Post
func New() *Session {
	return &Session{}
}

func (e *Session) MarshalBinary() (data []byte, err error) {
	sessionProto, err := Session2SessionProto(*e)
	if err != nil {
		return nil, err
	}
	return sessionProto.MarshalBinary()
}

func (e *Session) UnmarshalBinary(data []byte) (err error) {
	sessionProto := &proto.Session{}

	err = sessionProto.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	s, err := SessionProto2Session(*sessionProto)
	if err != nil {
		return err
	}

	*e = *s
	return nil
}
