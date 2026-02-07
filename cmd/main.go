package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const serverType = "TelemetryBridge"
const broadcastPort = 9999
const operationPort = 8888

type ServerInfo struct {
	Type string `json:"type"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type Broadcaster struct {
	Addr *net.UDPAddr
	Info ServerInfo
}

func (b *Broadcaster) Start(interval time.Duration) error {
	conn, err := net.DialUDP("udp", nil, b.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, _ := json.Marshal(b.Info)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		conn.Write(data)
	}
	return nil
}

func main() {
	// Get server IP (example: first non-loopback)
	addrs, _ := net.InterfaceAddrs()
	var serverIP string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			serverIP = ipnet.IP.String()
			break
		}
	}

	info := ServerInfo{
		Type: serverType,
		IP:   serverIP,
		Port: operationPort,
	}

	broadcastNet := "255.255.255.255"

	addr := &net.UDPAddr{IP: net.ParseIP(broadcastNet), Port: broadcastPort}
	b := Broadcaster{Addr: addr, Info: info}
	go func() {
		if err := b.Start(2 * time.Second); err != nil {
			os.Exit(1)
		}
	}()

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

		log.Printf("Received registration: %v", body)

		c.JSON(200, gin.H{"status": "OK", "data": body})
	})

	if err := router.Run(fmt.Sprintf(":%d", operationPort)); err != nil {
		os.Exit(1)
	}
}
