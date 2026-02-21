package digest

import "net"

const (
// serverType    = "TelemetryBridge"
// broadcastNet  = "255.255.255.255"
// broadcastPort = 9999
// operationPort = 8888
)

type Digest struct {
	Type          string `json:"type"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	PublicCertURL string `json:"public_cert_url"`
}

// func (b *Broadcaster) Start(interval time.Duration) error {
// 	conn, err := net.DialUDP("udp", nil, b.Addr)
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close()

// 	data, _ := json.Marshal(b.Info)
// 	ticker := time.NewTicker(interval)
// 	defer ticker.Stop()

//		for range ticker.C {
//			conn.Write(data)
//		}
//		return nil
//	}
func New() *Digest {
	addrs, _ := net.InterfaceAddrs()

	var publicIP string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			publicIP = ipnet.IP.String()
			break
		}
	}

	_ = publicIP

	return nil
}

// 	info := ServerInfo{
// 		Type: serverType,
// 		IP:   serverIP,
// 		Port: operationPort,
// 	}

// 	addr := &net.UDPAddr{IP: net.ParseIP(broadcastNet), Port: broadcastPort}
// 	b := Broadcaster{Addr: addr, Info: info}
// 	go func() {
// 		if err := b.Start(2 * time.Second); err != nil {
// 			os.Exit(1)
// 		}
// 	}()
