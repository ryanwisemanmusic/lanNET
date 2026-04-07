package discovery

import (
	"fmt"
	"net"
	"time"
)

const (
	broadcastPort = 12345
	discoveryMsg  = "LAN_GAME_DISCOVERY"
)

type GameServer struct {
	IP   string
	Port int
	Name string
}

func BroadcastPresence(serverName string, port int, stopChan <-chan bool) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("255.255.255.255:%d,", broadcastPort))
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	message := fmt.Sprintf("%s|%s|%d", discoveryMsg, serverName, port)

	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			conn.Write([]byte(message))
		}
	}
}

func DiscoverServers(timeout time.Duration) []GameServer {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", broadcastPort))
	if err != nil {
		return nil
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	servers := make(map[string]GameServer)
	buffer := make([]byte, 1024)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			break
		}
		msg := string(buffer[:n])
		var server GameServer
		if n, _ := fmt.Sscanf(msg, "%s|%s|%d", &server.IP, &server.Name, &server.Port); n == 3 {
			server.IP = remoteAddr.IP.String()
			servers[server.IP] = server
		}
	}

	result := make([]GameServer, 0, len(servers))
	for _, s := range servers {
		result = append(result, s)
	}
	return result
}
