package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	Id         uuid.UUID  `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	TotalPrice float64    `json:"total_price" gorm:"type:float8";not null;`
	UserID     uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty" gorm:"type:timestamp; default:NULL"`
	DeletedAt  time.Time  `json:"deleted_at,omitempty" gorm:"type:timestamp; default:NULL"`
	User       *Users     `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	CartItems  []CartItem `json:"cart_items,omitempty" gorm:"many2many:cart_item;"`
}

func (c *Cart) TableName() string {
	return "cart"
}

type CartItem struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	CartID    uuid.UUID `json:"cart_id" gorm:"type:uuid;not null;"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;"`
	Price     float64   `json:"price" gorm:"type:float8";not null;`
	Qty       int       `json:"qty" gorm:"type:int4;not null;default:1"`
	SubTotal  float64   `json:"subtotal" `
	CreatedAt time.Time `json:"created_at,omitempty"`
	Product   *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Cart      *Cart     `json:"cart,omitempty" gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;"`
}

func (c *CartItem) TableName() string {
	return "cart_item"
}
