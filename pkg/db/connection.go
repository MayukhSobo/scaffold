package db

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/pkg/log"
)

// Config holds database configuration
type Config struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	RetryAttempts   int           `mapstructure:"retry_attempts"`
	RetryDelay      time.Duration `mapstructure:"retry_delay"`
}

// NewConnection creates a new database connection using the provided configuration
func NewConnection(conf *viper.Viper, logger log.Logger) (*sql.DB, error) {
	config, err := parseConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	dsn := buildDSN(config)
	logger.Info("Connecting to database", log.String("host", config.Host), log.String("database", config.Name))

	db, err := connectWithRetry(dsn, config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", config.RetryAttempts, err)
	}

	// Configure connection pool
	configureConnectionPool(db, config)

	logger.Info("Database connection established successfully")
	return db, nil
}

// parseConfig extracts database configuration from Viper
func parseConfig(conf *viper.Viper) (*Config, error) {
	config := &Config{
		// Set defaults
		Host:            "localhost",
		Port:            "3306",
		User:            "root",
		Password:        "",
		Name:            "scaffold",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
		RetryAttempts:   5,
		RetryDelay:      2 * time.Second,
	}

	// Extract database configuration from db.mysql section
	if conf.IsSet("db.mysql") {
		if err := conf.UnmarshalKey("db.mysql", config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal db.mysql config: %w", err)
		}
		// Decode base64 password if needed
		config.Password = decodeIfBase64(config.Password)
	}

	// Override with individual keys if they exist
	if conf.IsSet("db.mysql.host") {
		config.Host = conf.GetString("db.mysql.host")
	}
	if conf.IsSet("db.mysql.port") {
		config.Port = conf.GetString("db.mysql.port")
	}
	if conf.IsSet("db.mysql.user") {
		config.User = conf.GetString("db.mysql.user")
	}
	if conf.IsSet("db.mysql.password") {
		config.Password = decodeIfBase64(conf.GetString("db.mysql.password"))
	}
	if conf.IsSet("db.mysql.database") {
		config.Name = conf.GetString("db.mysql.database")
	}

	// Also support legacy "database" key for backwards compatibility
	if conf.IsSet("database") {
		if err := conf.UnmarshalKey("database", config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal database config: %w", err)
		}
		// Decode base64 password if needed
		config.Password = decodeIfBase64(config.Password)
	}

	// Legacy database config overrides
	if conf.IsSet("database.host") {
		config.Host = conf.GetString("database.host")
	}
	if conf.IsSet("database.port") {
		config.Port = conf.GetString("database.port")
	}
	if conf.IsSet("database.user") {
		config.User = conf.GetString("database.user")
	}
	if conf.IsSet("database.password") {
		config.Password = decodeIfBase64(conf.GetString("database.password"))
	}
	if conf.IsSet("database.name") {
		config.Name = conf.GetString("database.name")
	}

	return config, nil
}

// buildDSN constructs the MySQL DSN string
func buildDSN(config *Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&tls=skip-verify&allowNativePasswords=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
}

// connectWithRetry attempts to connect to the database with retry logic
func connectWithRetry(dsn string, config *Config, logger log.Logger) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < config.RetryAttempts; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			logger.Warn("Failed to open database connection",
				log.Error(err),
				log.Int("attempt", i+1),
				log.Int("max_attempts", config.RetryAttempts),
			)
			time.Sleep(config.RetryDelay)
			continue
		}

		err = db.Ping()
		if err != nil {
			logger.Warn("Failed to ping database",
				log.Error(err),
				log.Int("attempt", i+1),
				log.Int("max_attempts", config.RetryAttempts),
			)
			db.Close()
			time.Sleep(config.RetryDelay)
			continue
		}

		// Connection successful
		logger.Info("Database connection established",
			log.Int("attempt", i+1),
		)
		return db, nil
	}

	return nil, err
}

// configureConnectionPool sets up the database connection pool parameters
func configureConnectionPool(db *sql.DB, config *Config) {
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
}

// MustConnect creates a database connection and panics on failure
// This is useful for application startup where database connectivity is critical
func MustConnect(conf *viper.Viper, logger log.Logger) *sql.DB {
	db, err := NewConnection(conf, logger)
	if err != nil {
		logger.Fatal("Critical: Unable to establish database connection", log.Error(err))
		panic(err) // This line won't be reached due to Fatal, but added for clarity
	}
	return db
}

// decodeIfBase64 decodes the password if it looks like base64 encoded data
func decodeIfBase64(value string) string {
	if value == "" {
		return value
	}

	// Simple check if it looks like base64 (alphanumeric + / + =)
	base64Pattern := regexp.MustCompile(`^[A-Za-z0-9+/]+=*$`)
	if !base64Pattern.MatchString(value) || len(value) <= 8 {
		return value
	}

	// Try to decode, if it fails, use original value
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return value
	}

	return string(decoded)
}
