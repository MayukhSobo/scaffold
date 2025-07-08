package main

import (
	"fmt"
	"os"

	"github.com/MayukhSobo/scaffold/pkg/config"
	"github.com/MayukhSobo/scaffold/pkg/db"
	"github.com/MayukhSobo/scaffold/pkg/log"
	"github.com/spf13/viper"
)

// Example: go run examples/config_environments.go
func configEnvironmentsExample() {
	fmt.Println("=== Configuration System Examples ===")
	fmt.Println("This demonstrates how to use configs with different environments")
	fmt.Println()

	// Example 1: Default configuration (local.yml)
	fmt.Println("1. Default Configuration (local.yml):")
	conf := config.NewConfig()
	showDatabaseConfig(conf, "local")
	fmt.Println()

	// Example 2: Docker configuration
	fmt.Println("2. Docker Configuration:")
	dockerConf := loadSpecificConfig("configs/docker.yml")
	if dockerConf != nil {
		showDatabaseConfig(dockerConf, "docker")
	}
	fmt.Println()

	// Example 3: Production configuration
	fmt.Println("3. Production Configuration:")
	prodConf := loadSpecificConfig("configs/prod.yml")
	if prodConf != nil {
		showDatabaseConfig(prodConf, "production")
	}
	fmt.Println()

	// Example 4: Configuration usage instructions
	fmt.Println("=== Usage Instructions ===")
	fmt.Println("To run your application with different configs:")
	fmt.Println("  Local:      ./server --config configs/local.yml")
	fmt.Println("  Docker:     ./server --config configs/docker.yml")
	fmt.Println("  Production: ./server --config configs/prod.yml")
	fmt.Println("  Alias:      ./server --config @/docker.yml")
	fmt.Println("  Validate:   ./server --config configs/local.yml --validate-config")
}

func loadSpecificConfig(configPath string) *viper.Viper {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("  Config file %s not found, skipping...\n", configPath)
		return nil
	}

	conf := viper.New()
	conf.SetConfigFile(configPath)
	if err := conf.ReadInConfig(); err != nil {
		fmt.Printf("  Error reading config %s: %v\n", configPath, err)
		return nil
	}
	return conf
}

func showDatabaseConfig(conf *viper.Viper, envType string) {
	// Environment info
	env := conf.GetString("env")
	fmt.Printf("  Environment: %s\n", env)

	// Database configuration
	dbHost := conf.GetString("db.mysql.host")
	dbPort := conf.GetString("db.mysql.port")
	dbUser := conf.GetString("db.mysql.user")
	dbName := conf.GetString("db.mysql.database")

	fmt.Printf("  Database Host: %s:%s\n", dbHost, dbPort)
	fmt.Printf("  Database: %s (user: %s)\n", dbName, dbUser)

	// Test database configuration parsing
	logger := log.NewConsoleLogger(log.InfoLevel)
	_, err := db.NewConnection(conf, logger)
	if err != nil {
		fmt.Printf("  ❌ Database config test failed: %v\n", err)
	} else {
		fmt.Printf("  ✅ Database config valid\n")
	}

	// Show other relevant settings
	httpPort := conf.GetInt("http.port")
	if httpPort > 0 {
		fmt.Printf("  HTTP Port: %d\n", httpPort)
	}

	logLevel := conf.GetString("log.level")
	if logLevel != "" {
		fmt.Printf("  Log Level: %s\n", logLevel)
	}
}
