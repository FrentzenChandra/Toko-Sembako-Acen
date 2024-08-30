package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	Id         uuid.UUID  `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	TotalPrice float64    `json:"total_price" gorm:"type:float8";not null;`
	Products	[]Product `json:"products,omitempty" `
	CreatedAt  *time.Time `json:"created_at,string,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at_at,string,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at_at,string,omitempty"`
}

func (e *Cart) TableName() string {
	return "carts"
}
