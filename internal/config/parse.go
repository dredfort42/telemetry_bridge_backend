package config

import (
	config "github.com/dredfort42/go_config_reader"
	"github.com/google/uuid"
)

func parseConfiguration(cfg *config.Config) {
	App.UUID = uuid.New().String()
	App.Debug = cfg.GetBool("general.debug")

	App.Service = ServiceConfig{
		Host:           cfg.GetString("server.host"),
		Port:           cfg.GetInt("server.port"),
		ReadTimeout:    cfg.GetDuration("server.read_timeout"),
		WriteTimeout:   cfg.GetDuration("server.write_timeout"),
		IdleTimeout:    cfg.GetDuration("server.idle_timeout"),
		MaxHeaderBytes: cfg.GetInt("server.max_header_bytes"),
	}

	App.Digester = DigesterConfig{
		BroadcastInterval: cfg.GetDuration("digester.broadcast_interval"),
		BroadcastAddress:  cfg.GetString("digester.broadcast_address"),
		BroadcastPort:     cfg.GetInt("digester.broadcast_port"),
	}
}
