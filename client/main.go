package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Connect error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Send a message: ")
		message, _ := reader.ReadString('\n')
		fmt.Fprint(conn, message)

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Server error: %v\n", err)
			return
		}
		fmt.Printf("Server says: %s", response)
	}
}
