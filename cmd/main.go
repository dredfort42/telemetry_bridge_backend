package main

import (
	"encoding/json"
	"net"
	"os"
	"time"
)

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
		Type: "TelemetryBridge",
		IP:   serverIP,
		Port: operationPort,
	}

	broadcastNet := "255.255.255.255"

	addr := &net.UDPAddr{IP: net.ParseIP(broadcastNet), Port: broadcastPort}
	b := Broadcaster{Addr: addr, Info: info}
	if err := b.Start(2 * time.Second); err != nil {
		os.Exit(1)
	}
}
