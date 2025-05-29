package main

import (
	"net"
	"strings"
)

type FilterConfig struct {
	Enabled         bool
	AllowedIPs      []string
	MinSeverity     int
	ExcludePatterns []string
}

type MessageFilter struct {
	config FilterConfig
	ipNets []*net.IPNet
}

func NewMessageFilter(config FilterConfig) (*MessageFilter, error) {
	filter := &MessageFilter{
		config: config,
		ipNets: make([]*net.IPNet, 0),
	}

	// Parse CIDR ranges from allowed IPs
	for _, ipStr := range config.AllowedIPs {
		if strings.Contains(ipStr, "/") {
			_, ipNet, err := net.ParseCIDR(ipStr)
			if err != nil {
				return nil, err
			}
			filter.ipNets = append(filter.ipNets, ipNet)
		}
	}

	return filter, nil
}

func (f *MessageFilter) ShouldProcessMessage(message string, addr net.Addr) bool {
	// If filtering is disabled, accept all messages
	if !f.config.Enabled {
		return true
	}

	// If no filters are configured, accept all messages
	if len(f.config.AllowedIPs) == 0 && f.config.MinSeverity == 7 && len(f.config.ExcludePatterns) == 0 {
		return true
	}

	// Check IP address
	if len(f.config.AllowedIPs) > 0 {
		ip := getIPFromAddr(addr)
		if !f.isIPAllowed(ip) {
			return false
		}
	}

	// Check severity
	severity := getSeverityFromMessage(message)
	if severity > f.config.MinSeverity {
		return false
	}

	// Check exclude patterns
	for _, pattern := range f.config.ExcludePatterns {
		if strings.Contains(message, pattern) {
			return false
		}
	}

	return true
}

func (f *MessageFilter) isIPAllowed(ip net.IP) bool {
	if ip == nil {
		return false
	}

	// Always allow localhost/loopback addresses
	if ip.IsLoopback() {
		return true
	}

	// Check exact IP matches
	for _, allowedIP := range f.config.AllowedIPs {
		if !strings.Contains(allowedIP, "/") {
			if allowedIP == ip.String() {
				return true
			}
		}
	}

	// Check CIDR ranges
	for _, ipNet := range f.ipNets {
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}

func getIPFromAddr(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.UDPAddr:
		return v.IP
	case *net.TCPAddr:
		return v.IP
	default:
		host, _, _ := net.SplitHostPort(addr.String())
		return net.ParseIP(host)
	}
}

func getSeverityFromMessage(message string) int {
	if len(message) > 0 && message[0] == '<' {
		end := strings.Index(message, ">")
		if end > 0 && end <= 4 {
			pri := message[1:end]
			// Extract severity (last 3 bits of PRI)
			if severity := strings.Split(pri, "")[0]; len(severity) > 0 {
				switch severity {
				case "0":
					return 0 // Emergency
				case "1":
					return 1 // Alert
				case "2":
					return 2 // Critical
				case "3":
					return 3 // Error
				case "4":
					return 4 // Warning
				case "5":
					return 5 // Notice
				case "6":
					return 6 // Informational
				case "7":
					return 7 // Debug
				}
			}
		}
	}
	return 7 // Default to Debug level if unable to parse
}
