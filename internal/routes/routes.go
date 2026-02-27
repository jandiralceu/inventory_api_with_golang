// Package routes wires HTTP handlers to Gin engine routes and applies
// global middleware such as CORS, tracing, and authentication.
package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jandiralceu/inventory_api_with_golang/docs" // imported so swagger can read embedded docs
	"github.com/jandiralceu/inventory_api_with_golang/internal/config"
	"github.com/jandiralceu/inventory_api_with_golang/internal/handlers"
	"github.com/jandiralceu/inventory_api_with_golang/internal/middleware"
	platform "github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// RouteConfig holds all handler dependencies required to register API routes.
type RouteConfig struct {
	RoleHandler *handlers.RoleHandler
}

// Setup creates a configured [gin.Engine] with global middleware, public and
// protected route groups, and the Swagger UI endpoint.
func Setup(routeConfig *RouteConfig, config *config.Config, jwtManager *platform.JWTManager) *gin.Engine {
	gin.ForceConsoleColor()

	if config.Env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.TraceIDMiddleware())
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
		roles := api.Group("/roles")
		{
			roles.GET("", routeConfig.RoleHandler.FindAllRoles)
			roles.GET("/:id", routeConfig.RoleHandler.FindRoleByID)
			roles.POST("", routeConfig.RoleHandler.CreateRole)
			roles.DELETE("/:id", routeConfig.RoleHandler.DeleteRole)
		}
	}

	// Swagger UI route.
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
