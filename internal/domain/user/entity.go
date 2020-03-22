package user

import (
	"strconv"
	"time"
)

const TableName = "user"

// User is the user entity
type User struct {
	ID				uint		`gorm:"PRIMARY_KEY" json:"id"`
	Name			string		`gorm:"type:varchar(100);UNIQUE;INDEX" json:"username"`
	Passhash		string		`gorm:"type:bytea" json:"-"`
	CreatedAt		time.Time
	UpdatedAt		time.Time
	DeletedAt		*time.Time	`gorm:"INDEX"`
}


func (e User) TableName() string {
	return TableName
}

// New func is a constructor for the User
func New() *User {
	return &User{}
}

func (e User) GetID() string {
	return strconv.Itoa(int(e.ID))
}


func (e User) GetName() string {
	return e.Name
}
