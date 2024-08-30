package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID   `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	Name        string      `json:"name" gorm:"type:varchar(255)";not null;`
	Stock       int64       `json:"stock" gorm:"type:int4";not null;default:0`
	Price       float64     `json:"price" gorm:"type:float8";not null;`
	Capital     float64     `json:"capital" gorm:"type:float8;not null"`
	CategoryIds []uuid.UUID `json:"category_ids gorm:"type:_uuid"`
	Picture     string      `json:"picture" gorm:"type:text;not null"`
	CreatedAt   *time.Time  `json:"created_at,string,omitempty"`
	UpdatedAt   *time.Time  `json:"updated_at_at,string,omitempty"`
	DeletedAt   *time.Time  `json:"deleted_at_at,string,omitempty"`
	Categorys   []Category  `json:"categorys,omitempty" gorm:"foreignKey:CategoryIds"`
}
