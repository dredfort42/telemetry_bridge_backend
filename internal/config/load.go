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
	// SSLEnabled                bool
	// SSLCertFile               string
	// SSLKeyFile                string
	// JWTAccessTokenExpiration  time.Duration
	// JWTRefreshTokenExpiration time.Duration
}

// // DatabaseConfig represents database configuration
// type DatabaseConfig struct {
// 	Driver         string
// 	Host           string
// 	Port           int
// 	Name           string
// 	User           string
// 	Password       string
// 	SSLEnabled     bool
// 	MigrationsPath string
// 	MaxOpenConn    int
// 	MaxIdleConn    int
// 	MaxLifetime    time.Duration
// 	MaxIdleTime    time.Duration
// }

// KafkaConfig is a struct for Kafka configuration
type KafkaConfig struct {
	Brokers           []string
	TopicSystemDigest string
	TopicSensorData   string
	ConsumerConfig    map[string]any
	ProducerConfig    map[string]any
}

// // RedisConfig is a struct for Redis configuration
// type RedisConfig struct {
// 	Host            string
// 	Port            int
// 	User            string
// 	Password        string
// 	DB              int
// 	PoolSize        int
// 	MinIdleConns    int
// 	MaxConnAge      time.Duration
// 	IdleTimeout     time.Duration
// 	PoolTimeout     time.Duration
// 	DialTimeout     time.Duration
// 	ReadTimeout     time.Duration
// 	WriteTimeout    time.Duration
// 	MaxRetries      int
// 	MaxRetryBackoff time.Duration
// }

type PublicURLs struct {
	Register string
	Data     string
}

// AppConfig represents complete application configuration
type AppConfig struct {
	UUID              string
	AppName           string
	Version           string
	Environment       string
	HeartbeatInterval time.Duration
	JWTSecret         string
	Debug             bool
	Service           ServiceConfig
	// Database          DatabaseConfig
	Kafka KafkaConfig
	// Redis             RedisConfig
	PublicURLs PublicURLs
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
		"debug":              false,
		"heartbeat_interval": 10 * time.Second,

		"service.host":             "localhost",
		"service.port":             8080,
		"service.read_timeout":     "30s",
		"service.write_timeout":    "30s",
		"service.idle_timeout":     "120s",
		"service.max_header_bytes": 1048576, // 1MB
		// "service.ssl_enabled":                  false,
		// "service.ssl_cert_file":                "",
		// "service.ssl_key_file":                 "",
		// "service.jwt_access_token_expiration":  "15m",
		// "service.jwt_refresh_token_expiration": "7d",

		// "database.driver":          "postgres",
		// "database.ssl_enabled":     false,
		// "database.migrations_path": "migrations",
		// "database.max_open_conn":   25,
		// "database.max_idle_conn":   5,
		// "database.max_lifetime":    "5m",
		// "database.max_idle_time":   "2m",

		"kafka.consumer_config": map[string]any{"auto.offset.reset": "latest"},
		"kafka.producer_config": map[string]any{"acks": "all"},

		// "redis.pool_size":         100,
		// "redis.min_idle_conns":    10,
		// "redis.max_conn_age":      "15m",
		// "redis.idle_timeout":      "5m",
		// "redis.pool_timeout":      "30s",
		// "redis.dial_timeout":      "5s",
		// "redis.read_timeout":      "5s",
		// "redis.write_timeout":     "5s",
		// "redis.max_retries":       3,
		// "redis.max_retry_backoff": "100ms",
	}

	// Try to load from config file with defaults and validation
	opts := &config.LoadOptions{
		DefaultValues: defaults,
		RequiredKeys: []string{
			"jwt_secret",

			"database.host",
			"database.port",
			"database.name",
			"database.user",
			"database.password",

			"kafka.brokers",
			"kafka.topic_system_digest",
			"kafka.topic_sensor_data",

			"redis.host",
			"redis.port",
			"redis.password",
			"redis.db",
		},
		ValidationFunc: func(data map[string]any) error {
			// Validate port range
			if port, ok := data["service.port"].(int); ok {
				if port < 1024 || port > 65535 {
					return fmt.Errorf("service port must be between 1 and 65535, got %d", port)
				}
			}

			// // Validate database port
			// if dbPort, ok := data["database.port"].(int); ok {
			// 	if dbPort < 1024 || dbPort > 65535 {
			// 		return fmt.Errorf("database port must be between 1 and 65535, got %d", dbPort)
			// 	}
			// }

			// // Validate Redis port
			// if redisPort, ok := data["redis.port"].(int); ok {
			// 	if redisPort < 1024 || redisPort > 65535 {
			// 		return fmt.Errorf("redis port must be between 1 and 65535, got %d", redisPort)
			// 	}
			// }

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
