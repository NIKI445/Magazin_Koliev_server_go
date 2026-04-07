package router

import (
	"RestApiGo/internal/config"
	"RestApiGo/internal/handlers"
	"RestApiGo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	cfg := config.Load()

	r := gin.Default()

	// CORS middleware
	r.Use(middleware.CORS(cfg))

	// Auth routes
	authHandler := handlers.NewAuthHandler()
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/signup", authHandler.Signup)
	}

	// API routes
	requestHandler := handlers.NewRequestHandler()
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/products", requestHandler.GetProducts)
		apiGroup.POST("/product", requestHandler.AddToCart)
		apiGroup.GET("/cart/:id", requestHandler.GetCart)
		apiGroup.POST("/cart/count/:id", requestHandler.UpdateCartCount)
		apiGroup.DELETE("/cart/:id", requestHandler.DeleteCartItem)
		apiGroup.DELETE("/cartAll/:id", requestHandler.ClearCart)
	}

	return r
}
