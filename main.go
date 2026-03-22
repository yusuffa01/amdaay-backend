package main

import (
	"log"
	"time"

	"amdaaybackend/controllers"
	"amdaaybackend/middlewares" // 👇 1. KITA PANGGIL ALAMAT POS SATPAMNYA DI SINI
	"amdaaybackend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. BACA KUNCI RAHASIA DARI .ENV DULU
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan nilai bawaan.")
	}

	// 2. SETELAH KUNCI TERBACA, BARU SAMBUNGKAN KE DATABASE
	models.ConnectDatabase()

	// 3. JALANKAN MESIN ROUTER
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 👇 2. SURUH SATPAM ANTI-SPAM BERJAGA DI PINTU DEPAN
	r.Use(middlewares.RateLimiter())

	r.Static("/uploads", "./uploads")
	r.GET("/api/menu", controllers.GetMenus)
	r.GET("/api/menu/:id", controllers.GetMenuByID)
	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)
	r.POST("/api/lupa-password", controllers.LupaPassword)
	r.POST("/api/reset-password", controllers.ResetPassword)
	r.GET("/api/profile", controllers.CekToken(), controllers.GetProfile)
	r.PUT("/api/profile", controllers.CekToken(), controllers.UpdateProfile)
	r.PUT("/api/profile/password", controllers.CekToken(), controllers.UbahPassword)
	r.POST("/api/menu", controllers.CekAdmin(), controllers.CreateMenu)
	r.POST("/api/menu/:id", controllers.CekAdmin(), controllers.UpdateMenu)
	r.DELETE("/api/menu/:id", controllers.CekAdmin(), controllers.DeleteMenu)
	r.PUT("/api/menu/:id/toggle", controllers.CekAdmin(), controllers.ToggleMenuStatus)
	r.GET("/api/users", controllers.CekAdmin(), controllers.GetUsers)
	r.PUT("/api/users/:id", controllers.CekAdmin(), controllers.UpdateUserByAdmin)
	r.DELETE("/api/users/:id", controllers.CekAdmin(), controllers.DeleteUser)

	r.Run(":8080")
}