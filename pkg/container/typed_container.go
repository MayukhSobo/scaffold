package container

import (
	"database/sql"

	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/internal/repository/users"
	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/log"
)

// TypedContainer provides type-safe dependency injection
// This version uses specific interfaces for better type safety
type TypedContainer struct {
	// Infrastructure
	config   *viper.Viper
	logger   log.Logger
	database *sql.DB

	// Repositories - Type-safe versions
	userRepository users.Querier
	// Add more repositories as interfaces are defined
	// productRepository products.Querier
	// orderRepository   orders.Querier

	// Services - Type-safe versions
	userService service.UserService
	// Add more services as interfaces are defined
	// productService service.ProductService
	// orderService   service.OrderService
}

// NewTypedContainer creates a new type-safe dependency container
func NewTypedContainer(config *viper.Viper, logger log.Logger, database *sql.DB) *TypedContainer {
	container := &TypedContainer{
		config:   config,
		logger:   logger,
		database: database,
	}

	// Initialize all dependencies
	container.initializeDependencies()

	return container
}

// initializeDependencies creates all repository and service instances
func (c *TypedContainer) initializeDependencies() {
	// Initialize repositories
	c.userRepository = users.New(c.database)

	// Initialize base service
	baseService := service.NewService(c.logger)

	// Initialize services with their dependencies
	c.userService = service.NewUserService(baseService, c.userRepository)

	// Future repositories and services can be added here
	// c.productRepository = products.New(c.database)
	// c.productService = service.NewProductService(baseService, c.productRepository)
}

// Infrastructure getters
func (c *TypedContainer) GetConfig() *viper.Viper {
	return c.config
}

func (c *TypedContainer) GetLogger() log.Logger {
	return c.logger
}

func (c *TypedContainer) GetDatabase() *sql.DB {
	return c.database
}

// Repository getters
func (c *TypedContainer) GetUserRepository() users.Querier {
	return c.userRepository
}

// Service getters
func (c *TypedContainer) GetUserService() service.UserService {
	return c.userService
}

// Future repository getters (example templates)
// func (c *TypedContainer) GetProductRepository() products.Querier {
//     return c.productRepository
// }

// func (c *TypedContainer) GetOrderRepository() orders.Querier {
//     return c.orderRepository
// }

// Future service getters (example templates)
// func (c *TypedContainer) GetProductService() service.ProductService {
//     return c.productService
// }

// func (c *TypedContainer) GetOrderService() service.OrderService {
//     return c.orderService
// }

// GetAllServices returns a struct containing all services for easy access
func (c *TypedContainer) GetAllServices() *AllServices {
	return &AllServices{
		User: c.userService,
		// Product: c.productService,
		// Order:   c.orderService,
	}
}

// AllServices provides a single struct containing all services
// This makes it easy to pass all services to routes/controllers
type AllServices struct {
	User service.UserService
	// Product service.ProductService
	// Order   service.OrderService
	// Email   service.EmailService
	// Auth    service.AuthService
}

// GetAllRepositories returns a struct containing all repositories for easy access
func (c *TypedContainer) GetAllRepositories() *AllRepositories {
	return &AllRepositories{
		User: c.userRepository,
		// Product: c.productRepository,
		// Order:   c.orderRepository,
	}
}

// AllRepositories provides a single struct containing all repositories
// This can be useful for testing or advanced scenarios
type AllRepositories struct {
	User users.Querier
	// Product products.Querier
	// Order   orders.Querier
}
