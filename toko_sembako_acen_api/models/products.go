package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id         uuid.UUID  `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	Name       string     `json:"name" gorm:"type:varchar(255)";not null;`
	Stock      int        `json:"stock" gorm:"type:int4";not null;default:0`
	Price      float64    `json:"price" gorm:"type:float8";not null;`
	Capital    float64    `json:"capital" gorm:"type:float8;not null"`
	Picture    string     `json:"picture" gorm:"type:text;"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	UpdatedAt  *time.Time  `json:"updated_at,omitempty" gorm:"type:timestamp; default:NULL"`
	DeletedAt  *time.Time  `json:"deleted_at,omitempty" gorm:"type:timestamp; default:NULL"`
	Categories []Category `json:"categories,omitempty" gorm:"many2many:product_category;"`
}

type ProductCategory struct {
	Id         uuid.UUID  `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	ProductID  uuid.UUID  `json:"product_id" gorm:"type:uuid;not null;"`
	CategoryID uuid.UUID  `json:"category_id" gorm:"type:uuid;not null;"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	Product    *Product   `json:"product,omitempty" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	Category   *Category  `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE;"`
}

func (p *ProductCategory) TableName() string {
	return "product_category"
}

func (p *Product) TableName() string {
	return "product"
}
