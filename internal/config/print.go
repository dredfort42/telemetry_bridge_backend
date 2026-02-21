package config

import (
	"log"
)

func printConfiguration(l *log.Logger) {
	l.Println()

	l.Println("========================================")
	l.Println()

	l.Printf("Application Configuration:\n")
	l.Printf("├─ Name: %s v%s\n", App.AppName, App.Version)
	l.Printf("├─ Environment: %s\n", App.Environment)
	l.Printf("├─ UUID: %s\n", App.UUID)
	l.Printf("├─ Heartbeat Interval: %v\n", App.HeartbeatInterval)
	l.Printf("├─ JWT Secret: %s\n", App.JWTSecret)
	l.Printf("└─ Debug Mode: %v\n", App.Debug)
	l.Println()

	l.Printf("Manager Configuration:\n")
	l.Printf("├─ Address: %s:%d\n", App.Service.Host, App.Service.Port)
	l.Printf("├─ Read Timeout: %v\n", App.Service.ReadTimeout)
	l.Printf("├─ Write Timeout: %v\n", App.Service.WriteTimeout)
	l.Printf("├─ Idle Timeout: %v\n", App.Service.IdleTimeout)
	l.Printf("├─ Max Header Bytes: %d\n", App.Service.MaxHeaderBytes)
	// l.Printf("├─ SSL Enabled: %v\n", App.Service.SSLEnabled)
	// if App.Service.SSLEnabled {
	// 	l.Printf("├─ SSL Cert: %s\n", App.Service.SSLCertFile)
	// 	l.Printf("├─ SSL Key: %s\n", App.Service.SSLKeyFile)
	// }
	// l.Printf("├─ JWT Access Token Expiration: %v\n", App.Service.JWTAccessTokenExpiration)
	// l.Printf("└─ JWT Refresh Token Expiration: %v\n", App.Service.JWTRefreshTokenExpiration)
	l.Println()

	// l.Printf("Database Configuration:\n")
	// l.Printf("├─ Driver: %s\n", App.Database.Driver)
	// l.Printf("├─ Address: %s:%d\n", App.Database.Host, App.Database.Port)
	// l.Printf("├─ Database: %s\n", App.Database.Name)
	// l.Printf("├─ User: %s\n", App.Database.User)
	// l.Printf("├─ SSL Enabled: %v\n", App.Database.SSLEnabled)
	// l.Printf("├─ Migrations Path: %s\n", App.Database.MigrationsPath)
	// l.Printf("├─ Connection Pool: %d/%d (max idle/open)\n", App.Database.MaxIdleConn, App.Database.MaxOpenConn)
	// l.Printf("├─ Max Lifetime: %v\n", App.Database.MaxLifetime)
	// l.Printf("└─ Max Idle Time: %v\n", App.Database.MaxIdleTime)
	// l.Println()

	l.Printf("Kafka Configuration:\n")
	l.Printf("├─ Brokers: %v\n", App.Kafka.Brokers)
	l.Printf("├─ Topic System Digest: %s\n", App.Kafka.TopicSystemDigest)
	l.Printf("├─ Topic Sensors Data: %s\n", App.Kafka.TopicSensorData)
	l.Printf("├─ Consumer Config: %v\n", App.Kafka.ConsumerConfig)
	l.Printf("└─ Producer Config: %v\n", App.Kafka.ProducerConfig)
	l.Println()

	// l.Printf("Redis Configuration:\n")
	// l.Printf("├─ Address: %s:%d\n", App.Redis.Host, App.Redis.Port)
	// l.Printf("├─ DB: %d\n", App.Redis.DB)
	// l.Printf("├─ Connection Pool: %d/%d (min idle/total)\n", App.Redis.MinIdleConns, App.Redis.PoolSize)
	// l.Printf("├─ Max Conn Age: %v\n", App.Redis.MaxConnAge)
	// l.Printf("├─ Idle Timeout: %v\n", App.Redis.IdleTimeout)
	// l.Printf("├─ Pool Timeout: %v\n", App.Redis.PoolTimeout)
	// l.Printf("├─ Dial Timeout: %v\n", App.Redis.DialTimeout)
	// l.Printf("├─ Read Timeout: %v\n", App.Redis.ReadTimeout)
	// l.Printf("├─ Write Timeout: %v\n", App.Redis.WriteTimeout)
	// l.Printf("├─ Max Retries: %d\n", App.Redis.MaxRetries)
	// l.Printf("└─ Max Retry Backoff: %v\n", App.Redis.MaxRetryBackoff)
	// l.Println()

	l.Printf("Public URLs Configuration:\n")
	l.Printf("├─ Register: %s\n", App.PublicURLs.Register)
	l.Printf("└─ Data: %s\n", App.PublicURLs.Data)
	l.Println()

	l.Println("========================================")
	l.Println()
}
