package session

import (
	"encoding/json"
	"time"

	"github.com/Kalinin-Andrey/redditclone/internal/domain/user"
)

const (
	TableName	= "session"
)

type Data map[string]interface{}

// Session is the session entity
type Session struct {
	ID				uint		`gorm:"PRIMARY_KEY" json:"id"`
	UserID			uint     	`sql:"type:int REFERENCES \"user\"(id)" json:"userId"`
	User			user.User	`gorm:"FOREIGNKEY:UserID" json:"author"`
	Json			string		`sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
	Data			Data		`gorm:"-"`

	CreatedAt		time.Time		`json:"created"`
	UpdatedAt		time.Time		`json:"updated"`
	DeletedAt		*time.Time		`gorm:"INDEX" json:"deleted"`
}


func (s Session) TableName() string {
	return TableName
}

// New func is a constructor for the Post
func New() *Session {
	return &Session{}
}

func (s *Session) SetDataByJson() error {

	if s.Json == "" || s.Json == "{}" {
		s.Data	= make(map[string]interface{}, 1)
		s.Json	= "{}"
		return nil
	}

	err := json.Unmarshal([]byte(s.Json), &s.Data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) SetJsonByData() error {

	if s.Data == nil {
		s.Data	= make(map[string]interface{}, 1)
		s.Json	= "{}"
		return nil
	}

	bytes, err := json.Marshal(s.Data)
	if err != nil {
		return err
	}
	s.Json = string(bytes)
	return nil
}