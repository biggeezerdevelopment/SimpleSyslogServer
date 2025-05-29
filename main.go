package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	// Set up logging
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	log.Printf("Starting SimpleSyslogServer...")

	// Load configuration
	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create message filter
	filterConfig := FilterConfig{
		Enabled:         config.Filter.Enabled,
		AllowedIPs:      config.Filter.AllowedIPs,
		MinSeverity:     config.Filter.MinSeverity,
		ExcludePatterns: config.Filter.ExcludePatterns,
	}
	filter, err := NewMessageFilter(filterConfig)
	if err != nil {
		log.Fatalf("Failed to create message filter: %v", err)
	}

	// Create error channels
	tcpErrors := make(chan error, 1)
	udpErrors := make(chan error, 1)

	// Start TCP server
	go func() {
		tcpServer := NewTCPServer(config.Server.Port, filter)
		tcpErrors <- tcpServer.Start()
	}()

	// Start UDP server
	go func() {
		udpServer := NewUDPServer(config.Server.Port, filter)
		udpErrors <- udpServer.Start()
	}()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for stop signal or error
	select {
	case err := <-tcpErrors:
		log.Printf("TCP server error: %v", err)
	case err := <-udpErrors:
		log.Printf("UDP server error: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	}

	log.Printf("Server shutting down...")
}
