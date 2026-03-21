package config

import (
	"context"
	"fmt"
	"time"

	config "github.com/dredfort42/go_config_reader"
	log "github.com/dredfort42/go_logger"
)

// ServiceConfig is a struct for streams server configuration
type ServiceConfig struct {
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

// DigesterConfig is a struct for digester configuration
type DigesterConfig struct {
	BroadcastInterval time.Duration
	BroadcastAddress  string
	BroadcastPort     int
}

// AppConfig represents complete application configuration
type AppConfig struct {
	UUID        string
	AppName     string
	Version     string
	Environment string
	Debug       bool
	Service     ServiceConfig
	Digester    DigesterConfig
}

// App is the global application configuration
var App AppConfig

func Init(ctx context.Context, configPath string) error {
	log.Info.Println("Initializing config...")

	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to create config instance: %v", err)
	}

	// Define comprehensive defaults
	defaults := map[string]any{
		"debug": false,

		"service.host":             "localhost",
		"service.port":             8080,
		"service.read_timeout":     "30s",
		"service.write_timeout":    "30s",
		"service.idle_timeout":     "120s",
		"service.max_header_bytes": 1048576, // 1MB

		"digester.broadcast_interval": "10s",
		"digester.broadcast_address":  "255.255.255.255",
		"digester.broadcast_port":     8888,
	}

	// Try to load from config file with defaults and validation
	opts := &config.LoadOptions{
		DefaultValues: defaults,
		RequiredKeys:  []string{},
		ValidationFunc: func(data map[string]any) error {
			// Validate service port range
			if port, ok := data["service.port"].(int); ok {
				if port < 1024 || port > 65535 {
					return fmt.Errorf("service port must be between 1 and 65535, got %d", port)
				}
			}

			// Validate digester port range
			if port, ok := data["digester.broadcast_port"].(int); ok {
				if port < 1024 || port > 65535 {
					return fmt.Errorf("digester broadcast port must be between 1 and 65535, got %d", port)
				}
			}

			return nil
		},
		IgnoreEnv: false,
	}

	var configFiles []string

	if configPath != "" {
		configFiles = append(configFiles, configPath)
	} else {
		configFiles = []string{"config.yaml", "config.yml", "config.json", "config.ini"}
	}

	loaded := false

	for _, file := range configFiles {
		err := cfg.LoadFromFile(file, opts)
		if err == nil {
			log.Info.Println("Loaded configuration from:", file)
			loaded = true
			break
		} else {
			log.Warning.Printf("Failed to load configuration from %s: %v", file, err)
		}
	}

	if !loaded {
		return fmt.Errorf("failed to load configuration from any of the specified files: %v", configFiles)
	}

	log.Info.Println("Configuration loaded successfully")

	parseConfiguration(cfg)

	if App.Debug {
		log.Debug.Println("Debug mode is enabled")
		printConfiguration(log.Debug)
	}

	return nil
}
