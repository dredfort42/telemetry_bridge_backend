package main

import (
	"fmt"
	"log"
	"os"
	"telemetry_bridge/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := digest.New(config, )

	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		body := make(map[string]any)
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid JSON"})
			return
		}

		log.Printf("Received registration: %v", body)

		c.JSON(200, gin.H{"status": "registered", "data": body})
	})

	router.POST("/data", func(c *gin.Context) {
		// 	  json["mac"] = WiFi.macAddress();
		// json["temperature_c"] = temperature;
		// json["humidity_percent"] = humidity;
		// json["timestamp"] = millis();
		body := make(map[string]any)
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid JSON"})
			return
		}

		log.Printf("Received data: %v", body)
		// log.Printf("Timestamp: %v", int64(body["timestamp"].(float64)))

		c.JSON(200, gin.H{"status": "OK", "data": body})
	})

	if err := router.Run(fmt.Sprintf(":%s", config.Server.Port)); err != nil {
		os.Exit(1)
	}
}
