# SimpleSyslogServer

A lightweight syslog server implementation in Go that supports both TCP and UDP protocols. This server can filter messages based on IP addresses, severity levels, and content patterns.

## Features

- Supports both TCP and UDP syslog protocols
- Configurable message filtering:
  - IP-based filtering with CIDR support
  - Severity level filtering
  - Content pattern exclusion
- Configurable output:
  - File logging
  - Console output
- YAML-based configuration

## Prerequisites

- Go 1.23 or higher
- Basic understanding of syslog protocol

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/SimpleSyslogServer.git
cd SimpleSyslogServer
```

2. Build the server:
```bash
go build -o syslog-server
```

## Configuration

The server is configured using `config.yaml`. Here's an example configuration:

```yaml
server:
  port: 514
  log_file: "syslog.log"  # Leave empty for no file logging
  console_output: true    # Set to false to disable console output

filter:
  enabled: False          # Set to false to disable all filtering
  allowed_ips:
    - "127.0.0.0/8"      # Localhost range
    - "192.168.0.0/24"   # Local network
    - "0.0.0.0/0"        # All IPs
  min_severity: 7        # 0=Emergency, 7=Debug
  exclude_patterns:
    - "test"
    - "debug"
```

### Configuration Options

#### Server Settings
- `port`: The port number to listen on (default: 514)
- `log_file`: Path to the log file (empty for no file logging)
- `console_output`: Enable/disable console output

#### Filter Settings
- `enabled`: Enable/disable all filtering
- `allowed_ips`: List of allowed IP addresses/ranges in CIDR notation
- `min_severity`: Minimum severity level to accept (0-7)
  - 0: Emergency
  - 1: Alert
  - 2: Critical
  - 3: Error
  - 4: Warning
  - 5: Notice
  - 6: Informational
  - 7: Debug
- `exclude_patterns`: List of patterns to exclude from logging

## Usage

1. Start the server:
```bash
./syslog-server
```

2. Configure your syslog clients to send logs to the server:
```bash
# Example using logger command
logger -n localhost -P 514 "Test message"
```

## Testing

You can test the server using the `logger` command or any syslog client:

```bash
# Test UDP
logger -n localhost -P 514 -d "Test UDP message"

# Test TCP
logger -n localhost -P 514 -T "Test TCP message"
```

## License

This project is licensed under the terms of the included LICENSE file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 