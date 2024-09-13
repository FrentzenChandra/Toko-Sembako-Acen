package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id             uuid.UUID `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:uuid;not null;"`
	TotalNetIncome *float64  `json:"total_net_income" gorm:"type:float8"`
	TotalPrice     *float64  `json:"total_price" gorm:"type:float8"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	User           *Users    `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func (c *Order) TableName() string {
	return "order"
}

type OrderItem struct {
	Id           uuid.UUID `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	AdminName    string    `json:"admin_name" gorm:"type:varchar(255);not null;"`
	ProductID    uuid.UUID `json:"product_id" gorm:"type:uuid;not null;"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null;"`
	OrderID      uuid.UUID `json:"order_id" gorm:"type:uuid;not null;"`
	SubNetIncome float64   `json:"sub_net_income" gorm:"type:float8;not null"`
	Sub_total    float64   `json:"sub_total" gorm:"type:float8;not null"`
	Qty          int       `json:"qty" gorm:"type:int4;not null;default:1"`
	Price        float64   `json:"price" gorm:"type:float8;not null"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	Order        *Order    `json:"order,omitempty" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;"`
	User         *Users    `json:"user,omitempty" gorm:"foreignKey:UserID;"`
	Product      *Product  `json:"user,omitempty" gorm:"foreignKey:ProductID;"`
}

func (c *OrderItem) TableName() string {
	return "order_item"
}
