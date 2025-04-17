package models

import (
	"time"
)

type AdminUser struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);unique" json:"username"`
	Password  string    `gorm:"type:varchar(255)" json:"-"` // Password is not exposed in JSON
	Email     string    `gorm:"type:varchar(100);unique" json:"email"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	CreatedAt time.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp" json:"updated_at"`
}

func (AdminUser) TableName() string {
	return "admin_users"
}
