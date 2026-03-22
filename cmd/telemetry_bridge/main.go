package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"telemetry_bridge/internal/config"
	"telemetry_bridge/internal/digester"

	log "github.com/dredfort42/go_logger"
)

// Version is set at build time using -ldflags. Default is "0.0.0" if not provided.
var Version = "0.0.0"

// Application constants.
const (
	AppName     = "SENSOTECH.net Telemetry Bridge"
	Environment = "development"
)

func main() {
	var (
		versionFlag = flag.Bool("version", false, "Show version information")
		helpFlag    = flag.Bool("help", false, "Show help information")
		configFlag  = flag.String("config", "", "Path to config file")
	)

	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s version %s\n", AppName, Version)
		os.Exit(0)
	}

	if *helpFlag {
		showUsage()
		os.Exit(0)
	}

	log.Info.Printf("%s version %s is starting...", AppName, Version)
	defer log.Info.Printf("%s version %s stopped successfully", AppName, Version)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	config.App.AppName = AppName
	config.App.Version = Version
	config.App.Environment = Environment

	if err := config.Init(ctx, *configFlag); err != nil {
		log.Error.Println("Config initialization failed:", err)
		showUsage()
		cancel()
	}

	var wg sync.WaitGroup
	d := digester.New(ctx, cancel)
	wg.Add(1)
	go func() {
		if err := d.Start(&wg); err != nil {
			log.Error.Println("Digester failed:", err)
			cancel()
		}
	}()

	wg.Wait()

	// if err := digest.New(config, )

	// router := gin.Default()
	// router.POST("/register", func(c *gin.Context) {
	// 	body := make(map[string]any)
	// 	if err := c.BindJSON(&body); err != nil {
	// 		c.JSON(400, gin.H{"error": "invalid JSON"})
	// 		return
	// 	}

	// 	log.Printf("Received registration: %v", body)

	// 	c.JSON(200, gin.H{"status": "registered", "data": body})
	// })

	// router.POST("/data", func(c *gin.Context) {
	// 	// 	  json["mac"] = WiFi.macAddress();
	// 	// json["temperature_c"] = temperature;
	// 	// json["humidity_percent"] = humidity;
	// 	// json["timestamp"] = millis();
	// 	body := make(map[string]any)
	// 	if err := c.BindJSON(&body); err != nil {
	// 		c.JSON(400, gin.H{"error": "invalid JSON"})
	// 		return
	// 	}

	// 	log.Printf("Received data: %v", body)
	// 	// log.Printf("Timestamp: %v", int64(body["timestamp"].(float64)))

	// 	c.JSON(200, gin.H{"status": "OK", "data": body})
	// })

	// if err := router.Run(fmt.Sprintf(":%s", config.Server.Port)); err != nil {
	// 	os.Exit(1)
	// }
}

// showUsage displays the usage information for the application.
func showUsage() {
	fmt.Printf(`%s - A telemetry bridge for IoT devices

Usage:
  %s [flags]

Flags:
  --version      Show version information and exit
  --help         Show this help message and exit
  --config PATH  Specify path to config file.
                Supports YAML, JSON and INI formats.

Example:
  %s --config /path/to/config.yaml

`, AppName, AppName, AppName)
}
