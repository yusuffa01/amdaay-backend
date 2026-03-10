package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"primaryKey;type:varchar(50)"` 
	Nama         string `gorm:"type:varchar(100);not null"`
	Email        string `gorm:"type:varchar(100);unique;not null"`
	Password     string `gorm:"type:varchar(255);not null"`
	Role         string `gorm:"type:varchar(20);default:'user'"`
	TokenReset   string `gorm:"type:varchar(100)"`
	ExpiredReset int64
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}