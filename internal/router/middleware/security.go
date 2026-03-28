package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Security creates security middleware for HTTP headers.
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// // Allow localhost and production origins
		// allowedOrigins := []string{
		// 	"http://localhost:8080",
		// 	"http://localhost:4200",
		// 	"http://localhost:4201",
		// 	"http://localhost:80",
		// 	"http://127.0.0.1:8080",
		// 	"http://127.0.0.1:4200",
		// 	"http://127.0.0.1:4201",
		// 	"http://127.0.0.1:80",
		// 	"http://media.srv.abm-jsc.ru",
		// 	"https://media.srv.abm-jsc.ru",
		// 	"https://v1.cavox.ru",
		// }

		// // Set CORS headers for allowed origins
		// if slices.Contains(allowedOrigins, origin) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		// }

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		// Only set HSTS for HTTPS
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
