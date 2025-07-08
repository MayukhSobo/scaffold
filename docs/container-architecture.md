# Container-Based Dependency Injection Architecture

## Overview

This document explains how the container pattern solves the scalability challenges when adding multiple services and repositories to your application. The container pattern provides a clean, maintainable, and testable approach to dependency injection.

## The Problem with Manual Dependency Injection

### Current Approach Issues

```go
// main.go becomes cluttered
userRepo := users.New(db)
productRepo := products.New(db)
orderRepo := orders.New(db)
paymentRepo := payments.New(db)

baseService := service.NewService(logger)
userService := service.NewUserService(baseService, userRepo)
productService := service.NewProductService(baseService, productRepo)
orderService := service.NewOrderService(baseService, orderRepo, userRepo, productRepo)
paymentService := service.NewPaymentService(baseService, paymentRepo, orderRepo)

// Every new service requires main.go changes
// Controllers need explicit injection for each dependency
// Testing becomes complex with multiple mocks
```

### Problems:
1. **main.go bloat**: Every new service requires changes to main.go
2. **Controller modifications**: Adding services requires changing controller constructors
3. **Tight coupling**: Controllers depend on specific service implementations
4. **Complex testing**: Mocking multiple dependencies is cumbersome
5. **Maintenance overhead**: Changes ripple through multiple files

## The Container Solution

### Core Concept

The container pattern centralizes all dependency management in a single location, providing clean access to any service or repository without tight coupling.

```go
// main.go stays clean regardless of number of services
database := db.MustConnect(conf, logger)
container := container.NewTypedContainer(conf, logger, database)

server.RunWithCustomSetup(conf, logger, func(s *server.FiberServer) {
    s.SetupBusinessRoutesWithContainer(container)
})
```

### Key Components

1. **TypedContainer**: Manages all dependencies with type safety
2. **Route Registration**: Uses container to access any service
3. **Handlers**: Get dependencies through container
4. **Testing**: Easy mocking via container

## Implementation Details

### 1. Container Structure

```go
type TypedContainer struct {
    // Infrastructure
    config   *viper.Viper
    logger   log.Logger
    database *sql.DB

    // Repositories - Type-safe versions
    userRepository    users.Querier
    productRepository products.Querier
    orderRepository   orders.Querier

    // Services - Type-safe versions
    userService    service.UserService
    productService service.ProductService
    orderService   service.OrderService
}
```

### 2. Dependency Initialization

```go
func (c *TypedContainer) initializeDependencies() {
    // Initialize repositories
    c.userRepository = users.New(c.database)
    c.productRepository = products.New(c.database)
    c.orderRepository = orders.New(c.database)

    // Initialize base service
    baseService := service.NewService(c.logger)

    // Initialize services with their dependencies
    c.userService = service.NewUserService(baseService, c.userRepository)
    c.productService = service.NewProductService(baseService, c.productRepository)
    c.orderService = service.NewOrderService(baseService, c.orderRepository, c.userRepository, c.productRepository)
}
```

### 3. Route Registration

```go
func RegisterOrderRoutesWithContainer(router fiber.Router, baseHandler *handler.Handler, container *container.TypedContainer) {
    // Container provides easy access to any service
    orderService := container.GetOrderService()
    userService := container.GetUserService()      // No additional injection needed
    productService := container.GetProductService() // No additional injection needed
    
    // Handler can access multiple services easily
    orderHandler := handler.NewOrderHandler(baseHandler, orderService, userService, productService)
    
    orders := router.Group("/orders")
    orders.Get("/", orderHandler.GetAllOrders)
    orders.Post("/", orderHandler.CreateOrder)
}
```

## Adding New Domains (Zero Main.go Changes)

### Step 1: Add to Container

```go
// In pkg/container/typed_container.go
type TypedContainer struct {
    // ... existing fields ...
    productRepository products.Querier      // Add this
    productService    service.ProductService // Add this
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

### Step 2: Create Route Registration

```go
// In internal/routes/product_routes.go
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

### Step 3: Register Routes

```go
// In internal/routes/routes_container.go
func RegisterRoutesWithContainer(crc *ContainerRouteConfig) {
    // ... existing code ...
    RegisterUserRoutesWithContainer(v1, baseHandler, crc.Container)
    RegisterProductRoutesWithContainer(v1, baseHandler, crc.Container)  // Add this line
}
```

**That's it!** No changes to main.go, server setup, or existing controllers.

## Benefits

### 1. Scalability
- Add unlimited services/repositories without modifying existing code
- Each domain is independent and self-contained
- No cascading changes when adding new features

### 2. Maintainability
- Single source of truth for all dependencies
- Clear separation of concerns
- Easy to refactor or replace components

### 3. Testability
```go
func TestOrderHandler(t *testing.T) {
    // Create mock container
    mockContainer := &container.TypedContainer{}
    
    // Inject mocks
    mockContainer.RegisterService("order", mockOrderService)
    mockContainer.RegisterService("user", mockUserService)
    
    // Test with container
    app := fiber.New()
    RegisterOrderRoutesWithContainer(app, baseHandler, mockContainer)
    
    // Run tests...
}
```

### 4. Type Safety
- Full compile-time checking
- IntelliSense support
- Interface-based design ensures contracts are met

### 5. Clean Architecture
- Loose coupling between components
- Dependency inversion principle
- Easy to understand and modify

## Migration Strategy

### Phase 1: Implement Container (Current)
- Create container package ✅
- Implement TypedContainer ✅
- Add container-based route registration ✅
- Keep existing approach working ✅

### Phase 2: Gradual Migration
- Add new domains using container pattern
- Migrate existing routes to container pattern
- Update tests to use container mocks

### Phase 3: Full Container Usage
- Remove manual dependency injection from main.go
- Standardize on container pattern throughout
- Update documentation and examples

## Real-World Example

Consider an e-commerce application with these domains:
- Users (authentication, profiles)
- Products (catalog, inventory)
- Orders (shopping cart, checkout)
- Payments (billing, transactions)
- Notifications (email, SMS)
- Analytics (tracking, reporting)

### Without Container
- main.go: 50+ lines of dependency setup
- Each new domain: 5-10 file modifications
- Testing: Complex mock setup for each test
- Maintenance: High coupling, difficult changes

### With Container
- main.go: 5 lines total
- Each new domain: 3 file additions, zero modifications
- Testing: Simple container mock injection
- Maintenance: Low coupling, easy changes

## Conclusion

The container pattern transforms your architecture from tightly coupled to loosely coupled, making it dramatically easier to scale and maintain as your application grows. It provides the foundation for clean, testable, and maintainable code that can handle any number of services and repositories without becoming unwieldy.

This approach is particularly valuable for:
- Microservice-like monoliths
- Applications with multiple business domains
- Teams with multiple developers
- Long-term maintenance and evolution
- Complex testing scenarios

The investment in setting up the container pattern pays dividends as soon as you add your second service domain, and the benefits compound with each additional domain. 