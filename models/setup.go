package models

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("amdaay.db"), &gorm.Config{})

	if err != nil {
		panic("Gagal terhubung ke database! Detail: " + err.Error())
	}

	database.AutoMigrate(&Menu{}, &User{}, &FotoProduk{})

	DB = database
}