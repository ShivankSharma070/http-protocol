package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("failed to resolve UDP address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("failed to dial UDP: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error reading input: %v", err)
			continue
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("error writing to UDP connection: %v", err)
		}
	}
}
