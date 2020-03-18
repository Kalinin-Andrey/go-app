package user

import (
	"strconv"
	"time"
)

const TableName = "user"

// Entity is the user entity
type Entity struct {
	ID				uint		`gorm:"PRIMARY_KEY" json:"id"`
	Name			string		`gorm:"type:varchar(100);UNIQUE;INDEX" json:"username"`
	//Salt			string		`gorm:"type:varchar(100)" json:"-"`
	Passhash		string		`gorm:"type:bytea" json:"-"`
	CreatedAt		time.Time
	UpdatedAt		time.Time
	DeletedAt		*time.Time	`gorm:"INDEX"`
}


func (e Entity) TableName() string {
	return TableName
}

// New func is a constructor for the Entity
func New() *Entity {
	return &Entity{}
}

func (e Entity) GetID() string {
	return strconv.Itoa(int(e.ID))
}


func (e Entity) GetName() string {
	return e.Name
}
