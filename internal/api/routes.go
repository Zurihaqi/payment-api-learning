package api

import (
	"payment-api-learning/internal/api/handlers"
	"payment-api-learning/internal/api/middleware"
	"payment-api-learning/internal/storage"

	"github.com/gin-gonic/gin"
)

func SetupRouter(storage storage.Storage) *gin.Engine {
	r := gin.Default()

	authHandler := handlers.NewAuthHandler(storage)
	paymentHandler := handlers.NewPaymentHandler(storage)

	r.POST("/login", authHandler.Login)

	authRoutes := r.Group("/")
	authRoutes.Use(middleware.AuthMiddleware(authHandler))
	{
		authRoutes.POST("/logout", authHandler.Logout)
		authRoutes.POST("/payment", paymentHandler.ProcessPayment)
	}

	return r
}
