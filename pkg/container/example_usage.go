package container

/*
EXAMPLE: How to scale your application with the Container Pattern

This file demonstrates how controllers handle multiple services and repositories
with minimal to no modifications when new domains are added.

Current architecture problems:
1. main.go gets cluttered with dependency creation
2. Controllers need manual dependency injection for each service
3. Adding new services requires changes throughout the codebase
4. Hard to test and mock dependencies

Container pattern solutions:
1. All dependencies managed in one place
2. Controllers access any service through the container
3. Adding new services requires minimal changes
4. Easy testing and mocking

## EXAMPLE 1: Current approach (manual dependency injection)

// In main.go - gets messy with more services
userRepo := users.New(db)
productRepo := products.New(db)
orderRepo := orders.New(db)
paymentRepo := payments.New(db)

baseService := service.NewService(logger)
userService := service.NewUserService(baseService, userRepo)
productService := service.NewProductService(baseService, productRepo)
orderService := service.NewOrderService(baseService, orderRepo, userRepo, productRepo)
paymentService := service.NewPaymentService(baseService, paymentRepo, orderRepo)

// In FiberServer.SetupBusinessRoutes - needs modification for each new service
func (s *FiberServer) SetupBusinessRoutes(
    userService service.UserService,
    productService service.ProductService,
    orderService service.OrderService,
    paymentService service.PaymentService,
) {
    // Route registration gets complex
}

## EXAMPLE 2: Container approach (scalable dependency injection)

// In main.go - stays clean regardless of number of services
database := db.MustConnect(conf, logger)
container := container.NewTypedContainer(conf, logger, database)

server.RunWithCustomSetup(conf, logger, func(s *server.FiberServer) {
    s.SetupBusinessRoutesWithContainer(container)
})

// In FiberServer.SetupBusinessRoutesWithContainer - no modification needed
func (s *FiberServer) SetupBusinessRoutesWithContainer(container *container.TypedContainer) {
    routeConfig := &routes.ContainerRouteConfig{
        App:       s.app,
        Container: container,
    }
    routes.RegisterRoutesWithContainer(routeConfig)
}

## EXAMPLE 3: Adding a new domain (Product) with ZERO main.go changes

Step 1: Add to TypedContainer
```go
type TypedContainer struct {
    // ... existing fields ...
    productRepository products.Querier  // Add this
    productService    service.ProductService  // Add this
}

func (c *TypedContainer) initializeDependencies() {
    // ... existing initialization ...
    c.productRepository = products.New(c.database)  // Add this
    c.productService = service.NewProductService(baseService, c.productRepository)  // Add this
}

func (c *TypedContainer) GetProductService() service.ProductService {  // Add this
    return c.productService
}
```

Step 2: Create product routes (new file)
```go
func RegisterProductRoutesWithContainer(router fiber.Router, baseHandler *handler.Handler, container *container.TypedContainer) {
    productService := container.GetProductService()
    productHandler := handler.NewProductHandler(baseHandler, productService)

    products := router.Group("/products")
    products.Get("/", productHandler.GetAllProducts)
    products.Get("/:id", productHandler.GetProductById)
    products.Post("/", productHandler.CreateProduct)
    products.Put("/:id", productHandler.UpdateProduct)
    products.Delete("/:id", productHandler.DeleteProduct)
}
```

Step 3: Add to route registration
```go
func RegisterRoutesWithContainer(crc *ContainerRouteConfig) {
    // ... existing code ...
    RegisterUserRoutesWithContainer(v1, baseHandler, crc.Container)
    RegisterProductRoutesWithContainer(v1, baseHandler, crc.Container)  // Add this line
}
```

That's it! No changes to main.go, server setup, or existing controllers.

## EXAMPLE 4: Complex service with multiple dependencies

```go
// Order service needs User, Product, and Payment services
func RegisterOrderRoutesWithContainer(router fiber.Router, baseHandler *handler.Handler, container *container.TypedContainer) {
    // Container provides easy access to any service
    orderService := container.GetOrderService()
    userService := container.GetUserService()      // No additional injection needed
    productService := container.GetProductService() // No additional injection needed

    // Handler can access multiple services easily
    orderHandler := handler.NewOrderHandler(baseHandler, orderService, userService, productService)

    orders := router.Group("/orders")
    orders.Get("/", orderHandler.GetAllOrders)
    orders.Get("/:id", orderHandler.GetOrderById)
    orders.Post("/", orderHandler.CreateOrder)
    orders.Put("/:id/status", orderHandler.UpdateOrderStatus)
}
```

## EXAMPLE 5: Testing becomes much easier

```go
func TestOrderHandler(t *testing.T) {
    // Create mock container
    mockContainer := &container.TypedContainer{}

    // Inject mocks
    mockContainer.RegisterService("order", mockOrderService)
    mockContainer.RegisterService("user", mockUserService)
    mockContainer.RegisterService("product", mockProductService)

    // Test with container
    app := fiber.New()
    baseHandler := handler.NewHandler(logger)
    RegisterOrderRoutesWithContainer(app, baseHandler, mockContainer)

    // Run tests...
}
```

## Benefits Summary:

1. **Scalability**: Add unlimited services/repositories without modifying existing code
2. **Maintainability**: Single source of truth for all dependencies
3. **Testability**: Easy mocking and dependency replacement
4. **Type Safety**: Full compile-time checking and IntelliSense
5. **Clean Architecture**: Clear separation of concerns
6. **Future-Proof**: Easy to refactor or replace components

The container pattern transforms your architecture from tightly coupled to loosely coupled,
making it much easier to scale and maintain as your application grows.
*/
