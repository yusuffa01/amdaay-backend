package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 1. Mengambil kunci rahasia dari brankas .env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 2. Merakit kunci menjadi satu kalimat akses
	// Penting: TiDB Serverless mewajibkan tambahan 'tls=true' di ujungnya
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=true", dbUser, dbPass, dbHost, dbPort, dbName)

	// 3. Membuka pintu brankas awan
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database TiDB! Detail: ", err)
	}

	// 4. Menyiapkan laci-laci tabel secara otomatis
	database.AutoMigrate(&Menu{}, &User{}, &FotoProduk{})

	DB = database
	fmt.Println("🚀 BERHASIL TERSAMBUNG KE BRANKAS BAJA TiDB!")
}