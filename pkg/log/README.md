# Log Package

Clean, pluggable logging package with dependency injection support.

## Architecture

This package follows proper dependency injection principles:
- **`logger.go`** - Core interface and field helpers
- **`consoleLogger.go`** - Console output implementation  
- **`fileLogger.go`** - File output with rotation implementation
- **`multiLogger.go`** - Combines multiple loggers

## Usage (Recommended Approach)

**Configuration-driven (Primary approach):**
```go
// ✅ RECOMMENDED: Use viper configuration for logger creation
import (
    "golang-di/pkg/config"
    "github.com/spf13/viper"
)

// Load configuration and create logger
conf := config.NewConfig()  // Loads from config files
logger := config.CreateLoggerFromConfig(conf)

// Use the logger
logger.Info("Application started", log.String("env", "production"))
```

**Configuration Format:**
```yaml
log:
  level: "debug"                    # debug, info, warn, error, fatal, panic
  console_logger:
    enabled: true
    colors: true                    # Enable colored output
    json_format: false              # Use structured JSON format
  file_logger:
    enabled: true
    directory: "logs"               # Directory for log files
    filename: "app.log"             # Log file name
    json_format: true               # JSON format for file logs
    max_size: 100                   # Max size in MB before rotation
    max_backups: 3                  # Number of backup files to keep
    max_age: 7                      # Days to keep old log files
    compress: true                  # Compress rotated files
```

**Direct creation (Testing/Special cases only):**
```go
// Only use direct constructors for testing or special cases
logger := log.NewConsoleLogger(log.InfoLevel)

fileLogger := log.NewFileLogger(log.InfoLevel, &log.FileLoggerConfig{
    Filename:   "app.log",
    MaxSize:    100,
    MaxBackups: 3,
    MaxAge:     7,
    Compress:   true,
    JsonFormat: true,
})

multiLogger := log.NewMultiLogger(consoleLogger, fileLogger)
```

## Structured Logging

```go
// Use field helpers for structured logging
logger.Info("User logged in",
    log.String("user_id", "123"),
    log.String("ip", "192.168.1.1"),
    log.Duration("response_time", time.Millisecond*150),
    log.Any("metadata", map[string]string{"browser": "chrome"}),
)

// Create logger with persistent context
contextLogger := logger.WithFields(
    log.String("request_id", "req-456"),
    log.String("component", "auth"),
)
contextLogger.Error("Authentication failed", log.Error(err))
```

## Benefits

- **Configuration-driven**: Logger behavior controlled by config files
- **Pluggable**: Easy to swap implementations
- **Framework-agnostic**: No external framework dependencies  
- **Clean interfaces**: Simple dependency injection
- **Type-safe fields**: Structured logging with helpers
- **Extensible**: Easy to add new logger types (DataDog, LogDNA, etc.)
- **Zero allocations**: Efficient logging with minimal overhead 