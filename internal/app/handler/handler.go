package handler

import (
	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/service"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests and manages routing.
// Contains dependencies required for request handlers (future fields).
type Handler struct {
	services *service.Service
}

// NewHandler creates and returns a new Handler instance.
// Constructor function for initializing a handler with possible dependencies.
func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}

}

// InitRoutes configures and returns the Gin router with defined endpoints.
// Adds middleware and registers handlers for all API paths.
func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	// Create a route group for subscription-related endpoints
	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.POST("/", h.createSubscription)                   //Create a new subscription
		subscriptions.GET("/", h.getAllSubscriptions)                   //Retrieve all subscriptions
		subscriptions.GET("/:subscription_id", h.getSubscriptionById)   //Get a specific subscription by ID
		subscriptions.PUT("/:subscription_id", h.updateSubscription)    //Update an existing subscription
		subscriptions.DELETE("/:subscription_id", h.deleteSubscription) //Delete a subscription
	}

	return router
}
