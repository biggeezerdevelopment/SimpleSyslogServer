package main

import (
	"fmt"
	"log"
	"net"
)

type UDPServer struct {
	port   int
	filter *MessageFilter
}

func NewUDPServer(port int, filter *MessageFilter) *UDPServer {
	return &UDPServer{
		port:   port,
		filter: filter,
	}
}

func (s *UDPServer) Start() error {
	addr := net.UDPAddr{
		Port: s.port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return fmt.Errorf("failed to start UDP server: %v", err)
	}
	defer conn.Close()

	log.Printf("UDP server listening on port %d", s.port)

	buffer := make([]byte, 65535) // Maximum UDP packet size

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading from UDP connection: %v", err)
			continue
		}

		message := string(buffer[:n])
		s.processSyslogMessage(message, remoteAddr)
	}
}

func (s *UDPServer) processSyslogMessage(message string, addr net.Addr) {
	if s.filter.ShouldProcessMessage(message, addr) {
		log.Printf("[UDP][%s] %s", addr.String(), message)
	}
}
