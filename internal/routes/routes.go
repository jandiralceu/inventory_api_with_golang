// Package routes wires HTTP handlers to Gin engine routes and applies
// global middleware such as CORS, tracing, and authentication.
package routes

import (
	"github.com/casbin/casbin/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jandiralceu/inventory_api_with_golang/docs" // imported so swagger can read embedded docs
	"github.com/jandiralceu/inventory_api_with_golang/internal/config"
	"github.com/jandiralceu/inventory_api_with_golang/internal/handlers"
	"github.com/jandiralceu/inventory_api_with_golang/internal/middleware"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// RouteConfig holds all handler dependencies required to register API routes.
type RouteConfig struct {
	AuthHandler      *handlers.AuthHandler
	RoleHandler      *handlers.RoleHandler
	UserHandler      *handlers.UserHandler
	CategoryHandler  *handlers.CategoryHandler
	SupplierHandler  *handlers.SupplierHandler
	WarehouseHandler *handlers.WarehouseHandler
	ProductHandler   *handlers.ProductHandler
	InventoryHandler *handlers.InventoryHandler
}

// Setup creates a configured [gin.Engine] with global middleware, public and
// protected route groups, and the Swagger UI endpoint.
func Setup(routeConfig *RouteConfig, config *config.Config, jwtManager *pkg.JWTManager, enforcer *casbin.Enforcer, cacheManager pkg.CacheManager) *gin.Engine {
	gin.ForceConsoleColor()

	router := gin.New()
	router.Use(middleware.TraceIDMiddleware())

	// Global Rate Limit per IP/User
	router.Use(middleware.RateLimiter(cacheManager, "global", config.RateLimitGlobal))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure CORS policy for cross-origin requests.
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:   []string{"Content-Length", "X-Trace-ID"},
		MaxAge:          12 * 3600,
	}))

	router.SetTrustedProxies(nil)

	router.Use(otelgin.Middleware(config.AppName))

	api := router.Group("/api/v1")
	{
		// Public routes (no authentication required).
		auth := api.Group("/auth")
		// Strict limit for auth
		auth.Use(middleware.RateLimiter(cacheManager, "auth", config.RateLimitAuth))
		{
			auth.POST("/signin", routeConfig.AuthHandler.SignIn)
			auth.POST("/register", routeConfig.AuthHandler.Register)
			auth.POST("/refresh", routeConfig.AuthHandler.RefreshToken)
			auth.POST("/signout", routeConfig.AuthHandler.SignOut)
		}

		// Protected routes (authentication required).
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		protected.Use(middleware.CasbinMiddleware(enforcer))
		{
			roles := protected.Group("/roles")
			{
				roles.GET("", routeConfig.RoleHandler.FindAllRoles)
				roles.GET("/:id", routeConfig.RoleHandler.FindRoleByID)
				roles.POST("", routeConfig.RoleHandler.CreateRole)
				roles.DELETE("/:id", routeConfig.RoleHandler.DeleteRole)
			}

			users := protected.Group("/users")
			{
				users.GET("", routeConfig.UserHandler.FindAllUsers)
				users.GET("/:id", routeConfig.UserHandler.FindUserByID)
				// Moderate limit for password changes
				users.PATCH("/change-password", middleware.RateLimiter(cacheManager, "pass-change", config.RateLimitAuth), routeConfig.UserHandler.ChangePassword)
				users.PATCH("/change-role", routeConfig.UserHandler.ChangeRole)
				users.DELETE("/:id", routeConfig.UserHandler.DeleteUser)
			}

			categories := protected.Group("/categories")
			{
				categories.GET("", routeConfig.CategoryHandler.FindAllCategories)
				categories.GET("/:id", routeConfig.CategoryHandler.FindCategoryByID)
				categories.POST("", routeConfig.CategoryHandler.CreateCategory)
				categories.PUT("/:id", routeConfig.CategoryHandler.UpdateCategory)
				categories.DELETE("/:id", routeConfig.CategoryHandler.DeleteCategory)
			}

			suppliers := protected.Group("/suppliers")
			{
				suppliers.GET("", routeConfig.SupplierHandler.FindAllSuppliers)
				suppliers.GET("/:id", routeConfig.SupplierHandler.FindSupplierByID)
				suppliers.POST("", routeConfig.SupplierHandler.CreateSupplier)
				suppliers.PUT("/:id", routeConfig.SupplierHandler.UpdateSupplier)
				suppliers.DELETE("/:id", routeConfig.SupplierHandler.DeleteSupplier)
			}

			warehouses := protected.Group("/warehouses")
			{
				warehouses.GET("", routeConfig.WarehouseHandler.FindAllWarehouses)
				warehouses.GET("/:id", routeConfig.WarehouseHandler.FindWarehouseByID)
				warehouses.POST("", routeConfig.WarehouseHandler.CreateWarehouse)
				warehouses.PUT("/:id", routeConfig.WarehouseHandler.UpdateWarehouse)
				warehouses.DELETE("/:id", routeConfig.WarehouseHandler.DeleteWarehouse)
			}

			products := protected.Group("/products")
			{
				products.GET("", routeConfig.ProductHandler.FindAllProducts)
				products.GET("/:id", routeConfig.ProductHandler.FindProductByID)
				products.POST("", routeConfig.ProductHandler.CreateProduct)
				products.PUT("/:id", routeConfig.ProductHandler.UpdateProduct)
				products.DELETE("/:id", routeConfig.ProductHandler.DeleteProduct)
			}

			inventory := protected.Group("/inventory")
			{
				inventory.GET("", routeConfig.InventoryHandler.FindAllInventory)
				inventory.GET("/transactions", routeConfig.InventoryHandler.GetTransactionHistory)
				inventory.GET("/:id", routeConfig.InventoryHandler.FindInventoryByID)
				inventory.POST("", routeConfig.InventoryHandler.CreateInventory)
				inventory.PUT("/:id", routeConfig.InventoryHandler.UpdateInventory)
				inventory.DELETE("/:id", routeConfig.InventoryHandler.DeleteInventory)

				// Stock operations
				inventory.POST("/:id/add", routeConfig.InventoryHandler.AddStock)
				inventory.POST("/:id/remove", routeConfig.InventoryHandler.RemoveStock)
				inventory.POST("/:id/reserve", routeConfig.InventoryHandler.ReserveStock)
				inventory.POST("/:id/release", routeConfig.InventoryHandler.ReleaseStock)
			}
		}
	}

	// Swagger UI route.
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
