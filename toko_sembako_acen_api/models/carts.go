package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Cart struct {
	Id         uuid.UUID  `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	TotalPrice float64    `json:"total_price" gorm:"type:float8";not null;`
	UserID     uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	User       *Users     `json:"user,omitempty" gorm:"foreignKey:UserID;"`
	CartItems  []CartItem `json:"cart_items,omitempty" gorm:"many2many:cart_item;constraint:OnDelete:CASCADE;""`
}

type CartItem struct {
	Id        uuid.UUID  `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	CartID    uuid.UUID  `json:"cart_id" gorm:"type:uuid;not null;"`
	ProductID uuid.UUID  `json:"product_id" gorm:"type:uuid;not null;"`
	Qty       int        `json:"qty" gorm:"type:int4;not null;default:1"`
	SubTotal  float64    `json:"subtotal" `
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Product   Product    `json:"product,omitempty" gorm:"foreignKey:ProductID;"`
	Cart      Cart       `json:"cart,omitempty" gorm:"foreignKey:CartID;"`
}

func (c *Cart) TableName() string {
	return "cart"
}

func (c *CartItem) TableName() string {
	return "cart_item"
}

func (u *Cart) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("cart_id = ?", u.Id).Delete(&CartItem{})
	return
}
