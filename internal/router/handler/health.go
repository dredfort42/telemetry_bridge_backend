package handler

import (
	"net/http"
	"runtime"
	"telemetry_bridge/internal/config"
	"time"

	"github.com/gin-gonic/gin"
)

// Health returns a health check response.
func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"ok":        true,
			"debug":     config.App.Debug,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"server": gin.H{
				"goroutines": runtime.NumGoroutine(),
				"mem_alloc":  mem.Alloc,
				"mem_total":  mem.TotalAlloc,
				"mem_sys":    mem.Sys,
			},
			"request": gin.H{
				"method": c.Request.Method,
				"path":   c.FullPath(),
				"remote": c.ClientIP(),
			},
		})
	}
}
