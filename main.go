package main

import (
	"time"
	"amdaaybackend/controllers"
	"amdaaybackend/models"
	"log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	models.ConnectDatabase()

err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan nilai bawaan.")
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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