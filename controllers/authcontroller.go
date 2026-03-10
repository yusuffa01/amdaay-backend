package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"amdaaybackend/models"
	"amdaaybackend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func getKunciRahasia() []byte {
	rahasia := os.Getenv("JWT_SECRET")
	if rahasia == "" {
		rahasia = "kunci_cadangan_sementara" 
	}
	return []byte(rahasia)
}

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	input.Password = string(hashedPassword)

	if err := models.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Registrasi berhasil, silakan login!"})
}

func Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email tidak ditemukan!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah!"})
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getKunciRahasia())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pesan": "Login berhasil!",
		"token": tokenString,
	})
}

func LupaPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format email salah"})
		return
	}

	var user models.User
	if err := models.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email tidak terdaftar!"})
		return
	}

	b := make([]byte, 16)
	rand.Read(b)
	tokenUnik := hex.EncodeToString(b)
	expiredAt := time.Now().Unix() + 300 

	models.DB.Model(&user).Updates(models.User{
		TokenReset:   tokenUnik,
		ExpiredReset: expiredAt,
	})

	resetLink := "http://localhost:5173/reset-password/" + tokenUnik
	err := utils.KirimEmailReset(user.Email, resetLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Link reset password telah dikirim ke email Anda!"})
}

func ResetPassword(c *gin.Context) {
	var input struct {
		Token        string `json:"token"`
		PasswordBaru string `json:"password_baru"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Data tidak lengkap"})
		return
	}

	var user models.User
	if err := models.DB.Where("token_reset = ?", input.Token).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "Link reset tidak valid atau sudah terpakai"})
		return
	}

	if time.Now().Unix() > user.ExpiredReset {
		c.JSON(400, gin.H{"error": "Link sudah kadaluwarsa! Silakan minta link baru."})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.PasswordBaru), bcrypt.DefaultCost)
	models.DB.Model(&user).Updates(map[string]interface{}{
		"Password":     string(hashedPassword),
		"TokenReset":   "",
		"ExpiredReset": 0,
	})

	c.JSON(200, gin.H{"message": "Password berhasil diperbarui!"})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	models.DB.Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"nama":  user.Nama,
		"email": user.Email,
	})
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input struct {
		Nama  string `json:"nama"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}
	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	models.DB.Model(&user).Updates(models.User{
		Nama:  input.Nama,
		Email: input.Email,
	})
	c.JSON(http.StatusOK, gin.H{"pesan": "Profil berhasil diperbarui!"})
}

func UbahPassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var input struct {
		PasswordLama string `json:"password_lama"`
		PasswordBaru string `json:"password_baru"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}
	var user models.User
	if err := models.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data user tidak ditemukan"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.PasswordLama)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password lama salah!"})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.PasswordBaru), bcrypt.DefaultCost)
	models.DB.Model(&user).Update("Password", string(hashedPassword))
	c.JSON(http.StatusOK, gin.H{"pesan": "Password berhasil diubah!"})
}


func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userIDLagiLogin, _ := c.Get("user_id")
	myID := fmt.Sprintf("%v", userIDLagiLogin)
	if myID == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda tidak bisa menghapus akun Anda sendiri!"})
		return
	}

	var user models.User
	if err := models.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	models.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"pesan": "Akun pengguna berhasil dihapus!"})
}

func UpdateUserByAdmin(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Nama     string `json:"nama"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	var user models.User
	if err := models.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	user.Nama = input.Nama
	user.Role = input.Role

	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	models.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"pesan": "Data pengguna berhasil diperbarui!"})
}

func CekToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Harap login terlebih dahulu!"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return getKunciRahasia(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal membaca token"})
			c.Abort()
		}
	}
}

func CekAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak!"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return getKunciRahasia(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			role, roleExists := claims["role"].(string)
			if !roleExists || role != "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Maaf, hanya Admin yang boleh mengakses ini!"})
				c.Abort()
				return
			}
            
			c.Set("user_id", claims["user_id"])
			c.Set("role", role)
			c.Next()
		}
	}
}