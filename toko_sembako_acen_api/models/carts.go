package models

import (
	"time"

	"github.com/google/uuid"
)

type CartItemInput struct {
	ProductID *uuid.UUID `json:"product_id"`
	Price     float64   `json:"price" `
	Qty       int       `json:"qty" `
}

type CartItem struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;"`
	Price     float64   `json:"price" gorm:"type:float8";not null;`
	Qty       int       `json:"qty" gorm:"type:int4;not null;default:1"`
	SubTotal  float64   `json:"subtotal" `
	CreatedAt time.Time `json:"created_at,omitempty"`
	Product   *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	User      *Users    `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func (c *CartItem) TableName() string {
	return "cart_item"
}
