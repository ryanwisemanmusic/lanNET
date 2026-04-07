package main

import (
	"bufio"
	"fmt"
	"lanNET/internal/network"
	"net"
	"os"
)

func main() {
	ip, err := network.GetLocalIPv4()
	if err != nil {
		fmt.Printf("Could not detect LAN IP: %v\n", err)
		os.Exit(1)
	}

	address := ip + ":8080"

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Listener error: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Server running on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client disconnected: %v\n", err)
			return
		}
		fmt.Fprintf(conn, "Echo: %s", message)
	}
}
