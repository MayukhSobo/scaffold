# Development Guide

Complete guide for development workflows, testing, and building with **Fiber web framework**, **MySQL database**, and **structured logging**.

## 🚀 Getting Started

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Task](https://taskfile.dev/installation/)
- [Docker & Docker Compose](https://www.docker.com/get-started)
- [MySQL](https://dev.mysql.com/downloads/) (for local development)
- [Python 3](https://www.python.org/downloads/) (for automation scripts)

### First-Time Setup

```bash
# Clone the repository
git clone https://github.com/MayukhSobo/scaffold.git
cd scaffold

# Install dependencies and setup environment
task deps:install

# Setup database (choose one)
task db:setup:local    # Local MySQL setup
task docker:compose:up # Docker environment
```

This will:
1. **Install Go dependencies** (`go mod download`)
2. **Install development tools** (golangci-lint, gotestsum, air, etc.)
3. **Setup database** with initial schema and data
4. **Generate SQLC code** for type-safe database operations
5. **Validate environment** setup

## 📊 Daily Development Workflow

### Starting Development

```bash
# Start development with hot reload (recommended)
task dev:hot

# Server will start on http://localhost:8000
# Hot reload with Air automatically restarts on code changes
```

### Database Development

```bash
# Start MySQL with Docker
task docker:compose:up

# Connect to database
task docker:compose:mysql

# Open Adminer (web database management)
task docker:compose:adminer

# Generate SQLC code after SQL changes
task db:generate:sqlc
```

### API Testing

```bash
# Test system endpoints
curl http://localhost:8000/health
curl http://localhost:8000/ping

# Test user API
curl http://localhost:8000/api/v1/users/admin
curl http://localhost:8000/api/v1/users/pending-verification
```

### Before Committing

```bash
# Quick checks (2-3 minutes)
task test:all
task quality:lint

# More thorough checks (5-10 minutes)  
task ci:pr
```

## 🧪 Testing

### Test Categories

```bash
# Run all tests
task test:all

# Specific test types
task test:unit          # Unit tests
task test:integration   # Integration tests
task test:race          # Race condition detection

# Database-specific tests
task test:db:all        # All database tests
task test:db:local      # Local MySQL tests (127.0.0.1)
task test:db:docker     # Docker MySQL tests (mysql hostname)

# Performance tests
task test:bench         # Benchmarks
```

### Coverage Analysis

```bash
# Generate and view coverage
task test:coverage:open

# Generate coverage treemap
task test:coverage:treemap
```

Coverage reports include:
- **Standard HTML reports**
- **Function-level analysis** with gocov
- **Visual treemaps** for large codebases

## 🏗️ Building

### Development Builds

```bash
# Debug build with race detection
task build:debug

# Run the debug binary
./build/scaffold-debug

# The application will start with:
# - Fiber server on configured port
# - Structured logging enabled
# - Database connection pooling
# - Hot reload middleware (if enabled)
```

### Release Builds

```bash
# Build for current platform
task build:release:linux    # Linux
task build:release:darwin   # macOS  
task build:release:windows  # Windows

# Build for all platforms
task build:release:all
```

Built binaries are located in:
- `build/linux/scaffold-amd64-linux`
- `build/darwin/scaffold-amd64-darwin`
- `build/windows/scaffold-amd64-windows.exe`

## 🐳 Docker Development

### Complete Environment

```bash
# Start MySQL, Adminer, and application
task docker:compose:up

# View logs from all services
task docker:compose:logs

# View logs from app only
task docker:compose:logs:app

# Restart all services
task docker:compose:restart

# Stop all services
task docker:compose:down
```

### Database Management

```bash
# Connect to MySQL shell
task docker:compose:mysql

# Open Adminer web interface
task docker:compose:adminer

# Generate MySQL init script
task db:generate:init
```

### Docker Builds

```bash
# Build application Docker image
task docker:build

# Run container
task docker:run

# Multi-platform builds
task docker:build:multi
```

## 🗄️ Database Development

### SQLC Workflow

1. **Write SQL queries** in `db/queries/*.sql`:
```sql
-- name: GetAdminUsers :many
SELECT * FROM users WHERE role = 'admin';

-- name: CreateUser :exec
INSERT INTO users (username, email, password_hash)
VALUES (?, ?, ?);
```

2. **Generate Go code**:
```bash
task db:generate:sqlc
```

3. **Use in repositories**:
```go
func (s *UserService) GetAdminUsers(ctx context.Context) ([]users.User, error) {
    return s.userRepo.GetAdminUsers(ctx)
}
```

### Database Testing

```bash
# Test database connections
task test:db:all

# Test specific environments
task test:db:local     # Tests local MySQL (127.0.0.1:3306)
task test:db:docker    # Tests Docker MySQL (mysql:3306)
```

### Migration Management

```bash
# Run migrations
task db:migrate

# Create new migration
task db:migration:create NAME=add_user_table

# Check migration status
task db:migration:status
```

## 📊 Logging & Monitoring

### Structured Logging

The application uses structured logging with smart field inclusion:

```json
{
  "level": "info",
  "time": "2024-01-15T10:30:45Z",
  "message": "HTTP Request",
  "method": "GET",
  "path": "/api/v1/users/admin",
  "status": 200,
  "latency": "15.67ms",
  "bytes_sent": "2.1KB"
}
```

### Log Configuration

```yaml
# configs/local.yml
log:
  level: "debug"
  loggers:
    console:
      driver: "console"
      enabled: true
      colors: true
    file:
      driver: "file"
      enabled: false
      directory: "logs"
      filename: "app.log"
```

### Viewing Logs

```bash
# Development logs (console)
task dev:hot

# Docker logs
task docker:compose:logs

# Application-specific logs
task docker:compose:logs:app -f  # Follow logs
```

## 🔍 Code Quality & Security

### Linting & Formatting

```bash
# Format code
task quality:fmt

# Run linter
task quality:lint

# Auto-fix issues
task quality:fix

# Run all quality checks
task quality:all
```

### Security Analysis

```bash
# Run security scanner
task quality:security

# Vulnerability scanning
task quality:gosec
```

## 🔧 Configuration Management

### Environment Configurations

| Environment | File | Description |
|-------------|------|-------------|
| **Local** | `configs/local.yml` | Development on localhost |
| **Docker** | `configs/docker.yml` | Docker Compose environment |
| **Production** | `configs/prod.yml` | Production settings |

### Configuration Structure

```yaml
env: local
app:
  name: "Scaffold v1.0.0"
  version: "1.0.0"
http:
  port: 8000
db:
  mysql:
    host: 127.0.0.1
    port: 3306
    user: scaffold
    password: my_secure_password_123
    database: user
server:
  middleware:
    logger: true
    cors: true
    recover: true
log:
  level: "debug"
  loggers:
    console:
      enabled: true
      colors: true
```

### Configuration Testing

```bash
# Validate configurations
task config:validate

# Test database connections
task config:test:db

# Show environment info
task config:info
```

## 🚀 Performance Optimization

### Benchmarking

```bash
# Run performance benchmarks
task test:bench

# Profile CPU usage
task profile:cpu

# Profile memory usage
task profile:mem
```

### Database Performance

```bash
# Test database connection pooling
task test:db:pool

# Analyze query performance
task db:analyze
```

### HTTP Performance

The Fiber framework provides excellent performance out of the box:
- **Built on Fasthttp** for maximum throughput
- **Zero-allocation routing** for low latency
- **Connection pooling** for database efficiency
- **Structured logging** with minimal overhead

## 🧩 Adding New Features

### Adding a New API Endpoint

1. **Create SQL queries** in `db/queries/`:
```sql
-- name: GetProductById :one
SELECT * FROM products WHERE id = ?;
```

2. **Generate SQLC code**:
```bash
task db:generate:sqlc
```

3. **Create service**:
```go
// internal/service/product.go
func (s *ProductService) GetProductById(ctx context.Context, id int64) (products.Product, error) {
    return s.productRepo.GetProductById(ctx, id)
}
```

4. **Create handler**:
```go
// internal/handler/product.go
func (h *ProductHandler) GetProductById(c *fiber.Ctx) error {
    id := c.ParamsInt("id")
    product, err := h.productService.GetProductById(c.Context(), int64(id))
    // ... handle response
}
```

5. **Register routes**:
```go
// internal/routes/product_routes.go
products.Get("/:id", productHandler.GetProductById)
```

### Using the Container Pattern

The application uses dependency injection for scalability:

```go
// Get services from container
userService := container.GetUserService()
productService := container.GetProductService()

// Container manages all dependencies
// No need to modify main.go for new services
```

## 🐛 Debugging

### Development Debugging

```bash
# Start with debug logging
task dev:hot

# Run with race detection
task build:debug
./build/scaffold-debug

# Enable verbose logging
LOG_LEVEL=debug task dev:hot
```

### Database Debugging

```bash
# Check database connection
task db:test

# View database logs
task docker:compose:logs mysql

# Connect to database shell
task docker:compose:mysql
```

### HTTP Debugging

```bash
# Test endpoints with curl
curl -v http://localhost:8000/health
curl -H "Accept: application/json" http://localhost:8000/api/v1/users/admin

# View request logs in real-time
task dev:hot  # Watch console output
```

## 🚨 Troubleshooting

### Common Issues

**Database Connection Failed:**
```bash
# Check MySQL is running
task docker:compose:up mysql

# Verify configuration
task config:test:db

# Check logs
task docker:compose:logs mysql
```

**Hot Reload Not Working:**
```bash
# Check Air configuration
cat .air.toml

# Restart Air
task dev:hot
```

**Tests Failing:**
```bash
# Run specific test category
task test:unit
task test:db:local

# Check test dependencies
task deps:install
```

### Getting Help

```bash
# Show all available tasks
task

# Show task help
task help

# Show environment information
task config:info
```

For detailed information about specific topics, see:
- [Container Architecture](container-architecture.md)
- [Task Reference](task-reference.md)
- [CI/CD Guide](ci-cd.md) 