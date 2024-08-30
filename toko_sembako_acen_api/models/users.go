package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Users struct {
	Id        uuid.UUID  `json:"id" gorm:"type:uuid;not null ; default:gen_random_uuid()"`
	Email     string     `json:"email" gorm:"type:varchar(255);not null"`
	Password  string     `json:"password" gorm:"type:text"`
	Username  string     `json:"username" gorm:"type:varchar(255);default:'steve'"`
	CreatedAt *time.Time `json:"created_at,string,omitempty"`
	UpdatedAt *time.Time `json:"updated_at_at,string,omitempty"`
	DeletedAt *time.Time `json:"deleted_at_at,string,omitempty"`
	Cart      Cart       `json:"cart,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

// TableName is Database TableName of this model

func (e *Users) TableName() string {
	return "users"
}

func (u *Users) AfterDelete(tx *gorm.DB) (err error) {
	// Ini berguna untuk Jika ada user 1 yang dihapus maka gorm akan mencari apakah ada
	// task yang memiliki user_id yang sama. Jika ada maka gorm akan menghapus task tersebut
	// juga secara otomatis
	tx.Clauses(clause.Returning{}).Where("user_id= ?", u.Id).Delete(&Cart{})
	return
}
