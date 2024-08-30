package models

import "time"

type Category struct {
	Id        int        `json:"id"`
	Name      string     `json:"name" gorm:"type:varchar(255)";not null;`
	CreatedAt *time.Time `json:"created_at,string,omitempty"`
	UpdatedAt *time.Time `json:"updated_at_at,string,omitempty"`

}


