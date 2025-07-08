# Testing Your Connected Routes

## Available Endpoints

After connecting the routes to the server, you now have these endpoints available:

### Basic Server Endpoints (Pre-existing)
- `GET /health` - Health check
- `GET /ping` - Ping endpoint  
- `GET /` - Root endpoint

### Your New Business Routes
- `GET /api/v1/users/admin` - Get users with admin access
- `GET /api/v1/users/pending-verification` - Get users with pending verification

## How to Test

### 1. Start the Server

#### With Default Configuration (local.yml)
```bash
# From project root
go run cmd/server/main.go
# or explicitly:
go run cmd/server/main.go --config configs/local.yml
```

#### With Different Environments
```bash
# Docker environment (port 12001)
go run cmd/server/main.go --config configs/docker.yml
go run cmd/server/main.go --config @/docker.yml

# Production environment (port 8080)  
go run cmd/server/main.go --config configs/prod.yml
go run cmd/server/main.go --config @/prod.yml

# Validate configuration
go run cmd/server/main.go --config configs/local.yml --validate-config
```

#### View Configuration Examples
```bash
# See how different environments load configs
go run examples/*.go config-examples
```

You should see logs showing:
```
Loaded config file: configs/local.yml
Starting application...
Initializing dependencies...
Database initialized
Repository layer initialized  
Service layer initialized
Starting server with business routes...
Business routes registered successfully
Server starting on port 8000
```

### 2. Test the Routes

#### Test Admin Users Route

**Note**: Port varies by environment:
- Local: `http://localhost:8000` (configs/local.yml)
- Docker: `http://localhost:12001` (configs/docker.yml)  
- Production: `http://localhost:8080` (configs/prod.yml)

```bash
curl http://localhost:8000/api/v1/users/admin
```

Expected Response:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "admin",
        "role": "admin",
        "status": "active"
      },
      {
        "id": 2,
        "username": "superadmin", 
        "role": "super_admin",
        "status": "active"
      }
    ],
    "count": 2
  }
}
```

#### Test Pending Verification Users Route
```bash
curl http://localhost:8000/api/v1/users/pending-verification
```

Expected Response:
```json
{
  "code": 0,
  "message": "success", 
  "data": {
    "users": [
      {
        "id": 3,
        "username": "user1",
        "email": "user1@example.com",
        "status": "pending_verification",
        "created_at": "2024-01-01T00:00:00Z",
        "verification_token": "abc123"
      },
      {
        "id": 4,
        "username": "user2",
        "email": "user2@example.com", 
        "status": "pending_verification",
        "created_at": "2024-01-02T00:00:00Z",
        "verification_token": "def456"
      }
    ],
    "count": 2
  }
}
```

#### Test Basic Endpoints
```bash
# Health check
curl http://localhost:8000/health

# Ping
curl http://localhost:8000/ping

# Root
curl http://localhost:8000/
```

## Architecture Overview

The connection flow is now:

```
main.go
├── Creates Database (mock)
├── Creates Repository Layer
├── Creates Service Layer  
├── Creates FiberServer
└── Calls SetupBusinessRoutes(userService)
    └── Registers /api/v1/users/* routes
```

## Next Steps

- **Database**: Replace `repository.NewDb()` with real database connection
- **Business Logic**: Replace mock data in handlers with real database queries
- **Authentication**: Add middleware for protected routes
- **Validation**: Add request validation for inputs 