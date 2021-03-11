package user

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	EntityName = "user"
	TableName  = "user"
)

// User is the user entity
type User struct {
	ID        uint   `gorm:"PRIMARY_KEY" json:"id"`
	Name      string `gorm:"type:varchar(100);UNIQUE;INDEX" json:"username"`
	Passhash  string `gorm:"type:bytea" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"INDEX"`
}

func (e User) TableName() string {
	return TableName
}

// New func is a constructor for the User
func New() *User {
	return &User{}
}

func (e User) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Name, validation.Required, validation.Length(2, 100), is.Alpha),
	)
}
