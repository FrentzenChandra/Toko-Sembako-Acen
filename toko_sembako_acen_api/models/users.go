package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id        uuid.UUID  `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	Email     string     `json:"email" gorm:"type:varchar(255);not null"`
	Password  string     `json:"password" gorm:"type:text"`
	Username  string     `json:"username" gorm:"type:varchar(255);default'steve'"`
	CreatedAt *time.Time `json:"created_at,string,omitempty"`
	UpdatedAt *time.Time `json:"updated_at_at,string,omitempty"`
	DeletedAt *time.Time `json:"deleted_at_at,string,omitempty"`
}

// TableName is Database TableName of this model

func (e *Users) TableName() string {
	return "users"
}

type OAuthUser struct {
	OAuthId       string `json:"oauthId"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
}
