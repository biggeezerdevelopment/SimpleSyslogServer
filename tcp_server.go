package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type TCPServer struct {
	port   int
	filter *MessageFilter
}

func NewTCPServer(port int, filter *MessageFilter) *TCPServer {
	return &TCPServer{
		port:   port,
		filter: filter,
	}
}

func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("TCP server listening on port %d", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting TCP connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				log.Printf("Error reading from TCP connection: %v", err)
			}
			return
		}

		s.processSyslogMessage(message, conn.RemoteAddr())
	}
}

func (s *TCPServer) processSyslogMessage(message string, addr net.Addr) {
	if s.filter.ShouldProcessMessage(message, addr) {
		log.Printf("[TCP][%s] %s", addr.String(), message)
	}
}
