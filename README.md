# 🚀 CodeEngine

> A containerized code execution service built with Go. Executes user-submitted code in isolated Docker environments with resource constraints and network restrictions.

## 📋 Overview

CodeEngine is a REST API service that accepts code snippets and executes them in secure, isolated Docker containers. Built to demonstrate containerization, security isolation, and concurrent request handling in Go.

## 🏗️ Architecture

```
CodeEngine/
├── cmd/                    # Application entry point
│   └── main.go            # HTTP server setup
├── internal/
│   ├── handlers/          # HTTP request handlers
│   ├── models/            # Data structures and types
│   └── services/          # Core business logic
└── Dockerfile             # Sandbox container definition
```

## ⚙️ Implementation Details

- 🔒 **Isolation**: Each execution runs in a fresh Docker container with no network access
- 📊 **Resource Control**: CPU (0.5 cores) and memory (50MB) limits prevent resource exhaustion
- 🔄 **Concurrency**: Semaphore-based concurrency control (max 10 simultaneous executions)
- ⏱️ **Timeout**: 30-second execution limit prevents infinite loops
- 🧹 **Cleanup**: Temporary files and containers are automatically removed

## 🛠️ Setup

### Prerequisites

- Go 1.23.3+
- Docker Engine
- Linux/macOS (recommended)

### Build and Run

```bash
# Clone and build
git clone <your-repo-url>
cd CodeEngine
go mod tidy
go build -o codeengine ./cmd

# Build sandbox container
docker build -t sandbox-python .

# Start server
./codeengine
```

Server runs on `http://localhost:8080`

## 📡 API

### Health Check
```http
GET /health
```

Returns: `OK`

### Code Execution
```http
POST /execute
Content-Type: application/json

{
  "language": "python3",
  "code": "print('Hello, World!')"
}
```

Response:
```json
{
  "output": "Hello, World!\n"
}
```

## 💡 Usage Examples

### Basic Execution
```bash
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "language": "python3",
    "code": "print(sum(range(10)))"
  }'
```

### Error Handling
```bash
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "language": "python3",
    "code": "print(undefined_variable)"
  }'
```

## ⚙️ Configuration

| Parameter | Value | Purpose |
|-----------|--------|---------|
| Max Concurrent | 10 | Limits simultaneous executions |
| Execution Timeout | 30s | Prevents runaway processes |
| CPU Limit | 0.5 cores | Resource constraint |
| Memory Limit | 50MB | Resource constraint |
| Network Access | None | Security isolation |

## 🔐 Security Model

- 🐳 **Container Isolation**: Fresh container per execution
- 🚫 **Network Disabled**: `--net=none` flag blocks all network access
- 📖 **Read-Only Code**: Source code mounted as read-only volume
- 🛡️ **Resource Limits**: CPU and memory constraints enforced by Docker
- 🗂️ **Temporary Filesystem**: All execution artifacts cleaned up automatically

## 🐍 Supported Languages

Currently supports Python 3.12. Architecture designed for easy extension to additional languages.

To add a new language:

1. Update language mapping in `models/types.go`
2. Create appropriate Dockerfile for the language runtime
3. Update executor to reference the new container image

## 🛠️ Development

```bash
# Run directly
go run ./cmd

# Run with hot reload (requires air)
go install github.com/cosmtrek/air@latest
air

# Test endpoints
curl http://localhost:8080/health
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"language": "python3", "code": "print(\"test\")"}'
```

## 🧠 Technical Considerations

This project demonstrates:

- 🌐 HTTP server implementation in Go
- 🐳 Docker API integration for container management
- ⚡ Concurrent request handling with semaphores
- 🔒 Security through containerization and resource limits
- 🏛️ Clean architecture with separation of concerns
- 📋 RESTful API design and JSON handling 