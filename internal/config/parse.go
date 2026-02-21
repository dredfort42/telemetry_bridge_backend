package config

import (
	config "github.com/dredfort42/go_config_reader"
	"github.com/google/uuid"
)

func parseConfiguration(cfg *config.Config) {
	App.Debug = cfg.GetBool("debug")
	App.HeartbeatInterval = cfg.GetDuration("heartbeat_interval")
	App.UUID = uuid.New().String()
	App.JWTSecret = cfg.GetString("jwt_secret")

	App.Service = ServiceConfig{
		Host:           cfg.GetString("manager.host"),
		Port:           cfg.GetInt("manager.port"),
		ReadTimeout:    cfg.GetDuration("manager.read_timeout"),
		WriteTimeout:   cfg.GetDuration("manager.write_timeout"),
		IdleTimeout:    cfg.GetDuration("manager.idle_timeout"),
		MaxHeaderBytes: cfg.GetInt("manager.max_header_bytes"),
		// SSLEnabled:                cfg.GetBool("manager.ssl_enabled"),
		// SSLCertFile:               cfg.GetString("manager.ssl_cert_file"),
		// SSLKeyFile:                cfg.GetString("manager.ssl_key_file"),
		// JWTAccessTokenExpiration:  cfg.GetDuration("manager.jwt_access_token_expiration"),
		// JWTRefreshTokenExpiration: cfg.GetDuration("manager.jwt_refresh_token_expiration"),
	}

	// App.Database = DatabaseConfig{
	// 	Driver:         cfg.GetString("database.driver"),
	// 	Host:           cfg.GetString("database.host"),
	// 	Port:           cfg.GetInt("database.port"),
	// 	Name:           cfg.GetString("database.name"),
	// 	User:           cfg.GetString("database.user"),
	// 	Password:       cfg.GetString("database.password"),
	// 	SSLEnabled:     cfg.GetBool("database.ssl_enabled"),
	// 	MigrationsPath: cfg.GetString("database.migrations_path"),
	// 	MaxOpenConn:    cfg.GetInt("database.max_open_conn"),
	// 	MaxIdleConn:    cfg.GetInt("database.max_idle_conn"),
	// 	MaxLifetime:    cfg.GetDuration("database.max_lifetime"),
	// 	MaxIdleTime:    cfg.GetDuration("database.max_idle_time"),
	// }

	App.Kafka = KafkaConfig{
		Brokers:           cfg.GetStringSlice("kafka.brokers"),
		TopicSystemDigest: cfg.GetString("kafka.topic_system_digest"),
		TopicSensorData:   cfg.GetString("kafka.topic_sensor_data"),
		ConsumerConfig:    cfg.GetNestedMap("kafka.consumer_config"),
		ProducerConfig:    cfg.GetNestedMap("kafka.producer_config"),
	}

	// App.Redis = RedisConfig{
	// 	Host:            cfg.GetString("redis.host"),
	// 	Port:            cfg.GetInt("redis.port"),
	// 	Password:        cfg.GetString("redis.password"),
	// 	DB:              cfg.GetInt("redis.db"),
	// 	PoolSize:        cfg.GetInt("redis.pool_size"),
	// 	MinIdleConns:    cfg.GetInt("redis.min_idle_conns"),
	// 	MaxConnAge:      cfg.GetDuration("redis.max_conn_age"),
	// 	IdleTimeout:     cfg.GetDuration("redis.idle_timeout"),
	// 	PoolTimeout:     cfg.GetDuration("redis.pool_timeout"),
	// 	DialTimeout:     cfg.GetDuration("redis.dial_timeout"),
	// 	ReadTimeout:     cfg.GetDuration("redis.read_timeout"),
	// 	WriteTimeout:    cfg.GetDuration("redis.write_timeout"),
	// 	MaxRetries:      cfg.GetInt("redis.max_retries"),
	// 	MaxRetryBackoff: cfg.GetDuration("redis.max_retry_backoff"),
	// }

	App.PublicURLs = PublicURLs{
		Register: cfg.GetString("public_urls.register"),
		Data:     cfg.GetString("public_urls.data"),
	}
}
