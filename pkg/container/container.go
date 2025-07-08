package container

import (
	"database/sql"

	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/pkg/log"
)

// Container holds all application dependencies in a centralized location
// This allows controllers to access any service without tight coupling
type Container struct {
	// Infrastructure
	config   *viper.Viper
	logger   log.Logger
	database *sql.DB

	// Repositories
	repositories *RepositoryContainer

	// Services
	services *ServiceContainer
}

// RepositoryContainer holds all repository instances
type RepositoryContainer struct {
	// Example repositories - add more as needed
	UserRepository    interface{} // This will be the specific repository interface
	ProductRepository interface{}
	OrderRepository   interface{}
	PaymentRepository interface{}
	// ... more repositories
}

// ServiceContainer holds all service instances
type ServiceContainer struct {
	// Example services - add more as needed
	UserService    interface{} // This will be the specific service interface
	ProductService interface{}
	OrderService   interface{}
	PaymentService interface{}
	EmailService   interface{}
	AuthService    interface{}
	// ... more services
}

// NewContainer creates a new dependency container
func NewContainer(config *viper.Viper, logger log.Logger, database *sql.DB) *Container {
	container := &Container{
		config:       config,
		logger:       logger,
		database:     database,
		repositories: &RepositoryContainer{},
		services:     &ServiceContainer{},
	}

	// Initialize repositories first
	container.initializeRepositories()

	// Initialize services (which depend on repositories)
	container.initializeServices()

	return container
}

// initializeRepositories creates all repository instances
func (c *Container) initializeRepositories() {
	// Initialize repositories here
	// Example: c.repositories.UserRepository = users.New(c.database)
	// This will be populated as we add more repositories
}

// initializeServices creates all service instances
func (c *Container) initializeServices() {
	// Initialize services here, injecting required repositories
	// Example: c.services.UserService = service.NewUserService(baseService, c.repositories.UserRepository)
	// This will be populated as we add more services
}

// Getters for infrastructure components
func (c *Container) GetConfig() *viper.Viper {
	return c.config
}

func (c *Container) GetLogger() log.Logger {
	return c.logger
}

func (c *Container) GetDatabase() *sql.DB {
	return c.database
}

// Getters for repositories
func (c *Container) GetRepositories() *RepositoryContainer {
	return c.repositories
}

// Getters for services
func (c *Container) GetServices() *ServiceContainer {
	return c.services
}

// GetUserRepository returns the user repository
func (r *RepositoryContainer) GetUserRepository() interface{} {
	return r.UserRepository
}

// GetProductRepository returns the product repository
func (r *RepositoryContainer) GetProductRepository() interface{} {
	return r.ProductRepository
}

// Add more repository getters as needed...

// GetUserService returns the user service
func (s *ServiceContainer) GetUserService() interface{} {
	return s.UserService
}

// GetProductService returns the product service
func (s *ServiceContainer) GetProductService() interface{} {
	return s.ProductService
}

// Add more service getters as needed...

// RegisterRepository allows dynamic registration of repositories
func (c *Container) RegisterRepository(name string, repository interface{}) {
	switch name {
	case "user":
		c.repositories.UserRepository = repository
	case "product":
		c.repositories.ProductRepository = repository
	case "order":
		c.repositories.OrderRepository = repository
	case "payment":
		c.repositories.PaymentRepository = repository
	default:
		c.logger.Warn("Unknown repository type for registration", log.String("name", name))
	}
}

// RegisterService allows dynamic registration of services
func (c *Container) RegisterService(name string, service interface{}) {
	switch name {
	case "user":
		c.services.UserService = service
	case "product":
		c.services.ProductService = service
	case "order":
		c.services.OrderService = service
	case "payment":
		c.services.PaymentService = service
	case "email":
		c.services.EmailService = service
	case "auth":
		c.services.AuthService = service
	default:
		c.logger.Warn("Unknown service type for registration", log.String("name", name))
	}
}
