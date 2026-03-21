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
	l.Printf("└─ Debug Mode: %v\n", App.Debug)
	l.Println()

	l.Printf("Service Configuration:\n")
	l.Printf("├─ Address: %s:%d\n", App.Service.Host, App.Service.Port)
	l.Printf("├─ Read Timeout: %v\n", App.Service.ReadTimeout)
	l.Printf("├─ Write Timeout: %v\n", App.Service.WriteTimeout)
	l.Printf("├─ Idle Timeout: %v\n", App.Service.IdleTimeout)
	l.Printf("└─ Max Header Bytes: %d\n", App.Service.MaxHeaderBytes)
	l.Println()

	l.Printf("Digester Configuration:\n")
	l.Printf("├─ Broadcast Interval: %v\n", App.Digester.BroadcastInterval)
	l.Printf("├─ Broadcast Address: %s\n", App.Digester.BroadcastAddress)
	l.Printf("└─ Broadcast Port: %d\n", App.Digester.BroadcastPort)
	l.Println()

	l.Println("========================================")
	l.Println()
}
