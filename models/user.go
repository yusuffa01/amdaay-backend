package models

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Nama         string `gorm:"type:varchar(100);not null"`
	Email        string `gorm:"type:varchar(100);unique;not null"`
	Password     string `gorm:"type:varchar(255);not null"`
	Role         string `gorm:"type:varchar(20);default:'user'"`
	TokenReset   string `gorm:"type:varchar(100)"`
	ExpiredReset int64
}