package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id        uuid.UUID  `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	Email     string     `json:"email" gorm:"type:varchar(255);not null"`
	Password  string     `json:"password" gorm:"type:text"`
	Username  string     `json:"username" gorm:"type:varchar(255);default:'steve'"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"type:timestamp; default:NULL"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"type:timestamp; default:NULL"`
}

// TableName is Database TableName of this model
func (u *Users) TableName() string {
	return "users"
}
