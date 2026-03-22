package digester

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"telemetry_bridge/internal/config"
	"time"

	log "github.com/dredfort42/go_logger"
)

// const (
// serverType    = "TelemetryBridge"
// broadcastNet  = "255.255.255.255"
// broadcastPort = 9999
// operationPort = 8888
// )

type Digester struct {
	Type       string             `json:"type"`
	IP         string             `json:"ip"`
	Port       int                `json:"port"`
	Ctx        context.Context    `json:"-"`
	CancelFunc context.CancelFunc `json:"-"`
}

func New(ctx context.Context, cancel context.CancelFunc) *Digester {
	addrs, _ := net.InterfaceAddrs()

	var LANIP string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			LANIP = ipnet.IP.String()
			break
		}
	}

	if config.App.Service.Host != "localhost" &&
		config.App.Service.Host != "127.0.0.1" &&
		config.App.Service.Host != "0.0.0.0" {
		LANIP = config.App.Service.Host
	}

	return &Digester{
		Type:       config.App.AppName,
		IP:         LANIP,
		Port:       config.App.Service.Port,
		Ctx:        ctx,
		CancelFunc: cancel,
	}
}

func (d *Digester) Start(wg *sync.WaitGroup) error {
	log.Info.Printf("Starting digester")       // Debug log for IP and port
	defer log.Info.Println("Digester stopped") // Ensure this is logged when the function exits
	defer wg.Done()

	conn, err := net.DialUDP(
		"udp",
		nil,
		&net.UDPAddr{
			IP:   net.ParseIP(config.App.Digester.BroadcastAddress),
			Port: config.App.Digester.BroadcastPort,
		},
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(d)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(config.App.Digester.BroadcastInterval)
	defer ticker.Stop()

	for {
		select {
		case <-d.Ctx.Done():
			return nil
		case <-ticker.C:
			conn.Write(data)
		}
	}
}
