package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MayukhSobo/scaffold/internal/repository"
	"github.com/MayukhSobo/scaffold/internal/server"
	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/config"
	"github.com/MayukhSobo/scaffold/pkg/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var (
	conf   *viper.Viper
	logger log.Logger
)

func init() {
	// Display startup banner
	fmt.Println(DisplayBanner())
	conf = config.NewConfig()
	var err error
	logger, err = log.CreateLoggerFromConfig(conf)
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}
}

func main() {
	logger.Info("Starting application...")

	// Create dependencies
	logger.Info("Initializing dependencies...")

	// Create database connection
	dbUser := conf.GetString("database.user")
	dbPassword := conf.GetString("database.password")
	dbHost := conf.GetString("database.host")
	dbPort := conf.GetString("database.port")
	dbName := conf.GetString("database.name")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	var db *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			logger.Error("failed to open database connection", log.Error(err))
			time.Sleep(2 * time.Second)
			continue
		}
		err = db.Ping()
		if err != nil {
			logger.Error("failed to ping database", log.Error(err))
			db.Close()
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		logger.Fatal("could not connect to the database after several retries", log.Error(err))
	}

	logger.Info("Database initialized")

	// Create repository layer
	queries := repository.New(db)
	userRepo := repository.NewUserRepository(queries)
	logger.Info("Repository layer initialized")

	// Create service layer
	baseService := service.NewService(logger)
	userService := service.NewUserService(baseService, userRepo)
	logger.Info("Service layer initialized")

	// Start server with custom setup to connect business routes
	logger.Info("Starting server with business routes...")
	server.RunWithCustomSetup(conf, logger, func(s *server.FiberServer) {
		// Setup business routes with dependencies
		s.SetupBusinessRoutes(userService)
		logger.Info("Business routes registered successfully")
	})
}
