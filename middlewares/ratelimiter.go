package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Catatan buku tamu untuk mencatat IP pelanggan
var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// Aturan: Maksimal 5 klik per detik. Kalau lebih, blokir!
		limiter = rate.NewLimiter(5, 10)
		visitors[ip] = limiter
	}
	return limiter
}

// RateLimiter adalah satpam penjaga pintu depan
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak permintaan! Sistem mendeteksi spam, silakan tunggu sebentar.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}