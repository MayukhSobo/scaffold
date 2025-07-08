package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/MayukhSobo/scaffold/internal/handler"
	"github.com/MayukhSobo/scaffold/pkg/container"
)

// ContainerRouteConfig holds the dependencies needed for route registration using container pattern
type ContainerRouteConfig struct {
	App       *fiber.App
	Container *container.TypedContainer
}

// RegisterRoutesWithContainer sets up all application routes using the container pattern
// This approach scales much better as you add more services and repositories
func RegisterRoutesWithContainer(crc *ContainerRouteConfig) {
	// Create base handler with logger from container
	baseHandler := handler.NewHandler(crc.Container.GetLogger())

	// Register API routes group
	api := crc.App.Group("/api")

	// Register v1 routes
	v1 := api.Group("/v1")

	// Register domain-specific routes
	RegisterUserRoutesWithContainer(v1, baseHandler, crc.Container)
	// Future route registrations - no modification needed to existing routes
	// RegisterProductRoutesWithContainer(v1, baseHandler, crc.Container)
	// RegisterOrderRoutesWithContainer(v1, baseHandler, crc.Container)
	// RegisterPaymentRoutesWithContainer(v1, baseHandler, crc.Container)
}

// RegisterUserRoutesWithContainer sets up user-related routes using container
func RegisterUserRoutesWithContainer(router fiber.Router, baseHandler *handler.Handler, container *container.TypedContainer) {
	// Get the user service from container
	userService := container.GetUserService()

	// Create user handler
	userHandler := handler.NewUserHandler(baseHandler, userService)

	// User routes group
	users := router.Group("/users")

	// Admin-specific user routes
	users.Get("/admin", userHandler.GetAdminUsers) // GET /api/v1/users/admin

	// Verification-specific user routes
	users.Get("/pending-verification", userHandler.GetPendingVerificationUsers) // GET /api/v1/users/pending-verification

	// Future user routes can be added here without affecting other modules
	// users.Get("/:id", userHandler.GetUserById)
	// users.Post("/", userHandler.CreateUser)
	// users.Put("/:id", userHandler.UpdateUser)
	// users.Delete("/:id", userHandler.DeleteUser)
}

// Example template for future route modules
// RegisterProductRoutesWithContainer sets up product-related routes using container
// func RegisterProductRoutesWithContainer(router fiber.Router, baseHandler *handler.Handler, container *container.TypedContainer) {
//     // Get the product service from container
//     productService := container.GetProductService()
//
//     // Create product handler
//     productHandler := handler.NewProductHandler(baseHandler, productService)
//
//     // Product routes group
//     products := router.Group("/products")
//
//     // Product-specific routes
//     products.Get("/", productHandler.GetAllProducts)
//     products.Get("/:id", productHandler.GetProductById)
//     products.Post("/", productHandler.CreateProduct)
//     products.Put("/:id", productHandler.UpdateProduct)
//     products.Delete("/:id", productHandler.DeleteProduct)
// }

// RegisterOrderRoutesWithContainer sets up order-related routes using container
// func RegisterOrderRoutesWithContainer(router fiber.Router, baseHandler *handler.Handler, container *container.TypedContainer) {
//     // Get multiple services from container if needed
//     orderService := container.GetOrderService()
//     userService := container.GetUserService()
//     productService := container.GetProductService()
//
//     // Create order handler with multiple service dependencies
//     orderHandler := handler.NewOrderHandler(baseHandler, orderService, userService, productService)
//
//     // Order routes group
//     orders := router.Group("/orders")
//
//     // Order-specific routes
//     orders.Get("/", orderHandler.GetAllOrders)
//     orders.Get("/:id", orderHandler.GetOrderById)
//     orders.Post("/", orderHandler.CreateOrder)
//     orders.Put("/:id/status", orderHandler.UpdateOrderStatus)
//     orders.Delete("/:id", orderHandler.CancelOrder)
// }

// RegisterAllRoutesWithContainer is a convenience function that registers all domain routes
// This demonstrates how the container pattern scales without modification
func RegisterAllRoutesWithContainer(crc *ContainerRouteConfig) {
	// Create base handler
	baseHandler := handler.NewHandler(crc.Container.GetLogger())

	// Register API routes group
	api := crc.App.Group("/api")
	v1 := api.Group("/v1")

	// Register all domain routes - each is independent and scalable
	RegisterUserRoutesWithContainer(v1, baseHandler, crc.Container)
	// Uncomment as you implement these modules:
	// RegisterProductRoutesWithContainer(v1, baseHandler, crc.Container)
	// RegisterOrderRoutesWithContainer(v1, baseHandler, crc.Container)
	// RegisterPaymentRoutesWithContainer(v1, baseHandler, crc.Container)
	// RegisterNotificationRoutesWithContainer(v1, baseHandler, crc.Container)
	// RegisterAnalyticsRoutesWithContainer(v1, baseHandler, crc.Container)
}
