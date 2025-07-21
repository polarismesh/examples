# Go-Zero + Polaris Service Discovery Example

This project demonstrates how to build distributed services using the go-zero microservice framework combined with Polaris service discovery. The project includes a gRPC server and a client, showcasing the complete workflow of service registration, discovery, and invocation.

## Project Structure

```
.
├── server/                 # gRPC Server
│   ├── server.go          # Server main program
│   ├── go.mod             # Go module dependencies
│   ├── go.sum             # Dependency lock file
│   ├── polaris.yaml       # Polaris configuration file
│   ├── etc/
│   │   └── polaris.yaml   # Server configuration
│   └── internal/
│       └── config/
│           └── config.go  # Configuration structure definition
├── client/                 # gRPC Client
│   ├── client.go          # Client main program
│   ├── go.mod             # Go module dependencies
│   ├── go.sum             # Dependency lock file
│   ├── polaris.yaml       # Polaris configuration file
│   └── etc/
│       └── polaris.yaml   # Client configuration
└── README.md              # Project documentation
```

## Environment Requirements

- **Go Version**: 1.18 or higher
- **Polaris Server**: Requires running Polaris service discovery service
- **Operating System**: Supports Linux, macOS, Windows

## Dependencies

### Main Dependencies
- [go-zero](https://github.com/zeromicro/go-zero): Microservice framework
- [zero-contrib](https://github.com/zeromicro/zero-contrib): go-zero extension components
- [Polaris](https://polarismesh.cn/): Service discovery and governance platform

### Core Features
- **Service Registration**: Server automatically registers with Polaris
- **Service Discovery**: Client discovers services through Polaris
- **Load Balancing**: Automatic client-side load balancing
- **Health Checking**: Service health status monitoring

## Quick Start

### 1. Start Polaris Service

First, you need to start the Polaris service discovery service. You can use one of the following methods:

#### Using Docker (Recommended)
```bash
# Pull and start Polaris service
docker run -d --name polaris -p 8091:8091 -p 8090:8090 polarismesh/polaris-server:latest
```

#### Using Docker Compose
```bash
# Create docker-compose.yml file
cat > docker-compose.yml << EOF
version: '3.8'
services:
  polaris:
    image: polarismesh/polaris-server:latest
    ports:
      - "8091:8091"
      - "8090:8090"
    environment:
      - POLARIS_LOG_LEVEL=info
EOF

# Start service
docker-compose up -d
```

### 2. Build and Run Server

```bash
# Enter server directory
cd server

# Download dependencies
go mod tidy

# Build server
go build -mod=vendor -o server server.go

# Run server
./server -f ./etc/polaris.yaml
```

After the server starts, it will:
- Listen on `0.0.0.0:3453` port
- Automatically register with Polaris service discovery center
- Service name is `EchoServerZero`

### 3. Build and Run Client

```bash
# Open a new terminal, enter client directory
cd client

# Download dependencies
go mod tidy

# Build client
go build -mod=vendor -o client client.go

# Run client
./client -f ./etc/polaris.yaml
```

After the client starts, it will:
- Discover `EchoServerZero` service through Polaris
- Send requests to the server every second
- Print server response results

## Configuration

### Server Configuration (`server/etc/polaris.yaml`)
```yaml
Name: server              # Service name
ListenOn: 0.0.0.0:3453   # Listen address and port
```

### Client Configuration (`client/etc/polaris.yaml`)
```yaml
Name: client                                                    # Client name
Target: polaris://127.0.0.1:38091/EchoServerZero?timeout=5000s # Target service address
```

### Polaris Connection Configuration (`polaris.yaml`)
```yaml
global:
  serverConnector:
    addresses:
      - 127.0.0.1:8091    # Polaris service address
```

## Testing Guide

### 1. Basic Functionality Test

After starting both server and client, the client should be able to successfully call the server:

```bash
# Client output example
15:04:05 => hello from hostname-xxx
15:04:06 => hello from hostname-xxx
15:04:07 => hello from hostname-xxx
```

### 2. Multi-Instance Test

Start multiple server instances to test load balancing:

```bash
# Terminal 1: Start first server instance
cd server
./server -f ./etc/polaris.yaml

# Terminal 2: Modify port and start second instance
cd server
# Temporarily modify port in config file to 3454
./server -f ./etc/polaris.yaml

# Terminal 3: Start client
cd client
./client -f ./etc/polaris.yaml
```

Observe the client output, you should see requests being distributed between different server instances.

### 3. Service Discovery Test

Test dynamic service up/down:

1. Start client and server
2. Stop server process
3. Observe if client can detect service unavailability
4. Restart server
5. Observe if client can automatically reconnect

### 4. Performance Test

You can use the following commands for simple performance testing:

```bash
# Use go test for benchmark testing
cd server
go test -bench=. -benchmem

# Or use hey tool for stress testing
# First install hey: go install github.com/rakyll/hey@latest
hey -n 1000 -c 10 -m POST -H "Content-Type: application/grpc" http://localhost:3453
```

## Troubleshooting

### Common Issues

1. **Failed to connect to Polaris**
   ```
   Error message: connection refused
   Solution: Ensure Polaris service is running, check if port 8091 is accessible
   ```

2. **Service registration failed**
   ```
   Error message: register service failed
   Solution: Check if the service address in polaris.yaml configuration file is correct
   ```

3. **Client failed to discover service**
   ```
   Error message: service not found
   Solution: Ensure server is started and successfully registered, check if service name matches
   ```

4. **Port conflict**
   ```
   Error message: bind: address already in use
   Solution: Modify port number in configuration file or stop process occupying the port
   ```

### Debugging Tips

1. **View detailed logs**
   ```bash
   # Add verbose logging when starting
   ./server -f ./etc/polaris.yaml -v
   ```

2. **Check Polaris console**
   - Visit `http://localhost:8090` to view Polaris management interface
   - Check service registration status and health check results

3. **Network connection test**
   ```bash
   # Test Polaris connection
   telnet 127.0.0.1 8091
   
   # Test server connection
   telnet 127.0.0.1 3453
   ```

## Development Guide

### Adding New gRPC Methods

1. Modify proto file definition
2. Regenerate Go code
3. Implement new method in server
4. Add invocation code in client

### Custom Configuration

You can customize the following by modifying configuration files:
- Service listening port
- Polaris service address
- Timeout duration
- Retry strategy

### Extension Features

The project supports the following extensions:
- Add middleware (authentication, rate limiting, monitoring)
- Integrate distributed tracing
- Add metrics monitoring
- Implement circuit breaker and fallback

## Related Links

- [go-zero Official Documentation](https://go-zero.dev/)
- [Polaris Official Documentation](https://polarismesh.cn/docs/)
- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/)
- [Project Source Code](https://github.com/polarismesh/examples)

## License

This project is licensed under the MIT License. See the LICENSE file for details. 