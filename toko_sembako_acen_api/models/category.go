package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	Id        uuid.UUID  `json:"id" gorm:"type:uuid;not null;default:gen_random_uuid()"`
	Name      string     `json:"name" gorm:"type:varchar(255);unique;not null;"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}


func (c *Category) TableName() string {
	return "category"
}
