package handler

import (
	"net/http"

	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/models"
	"github.com/gin-gonic/gin"
)

// createSubscription handles HTTP POST request for creating a new subscription
// This is the API endpoint handler for subscription creation
func (h *Handler) createSubscription(c *gin.Context) {
	// Bind JSON request body to Subscription model
	// Validates required fields based on 'binding' tags in the model
	var sub models.Subscription
	if err := c.BindJSON(&sub); err != nil {
		// Return 400 Bad Request if JSON is malformed or validation fails
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Call service layer to handle business logic and persistence
	// Service will validate business rules and create the subscription
	subId, err := h.services.SubscriptionStore.Create(sub)
	if err != nil {
		// Return 500 Internal Server Error if service operation fails
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Return 200 OK with the ID of the newly created subscription
	// The ID can be used by clients for subsequent operations
	c.JSON(http.StatusOK, map[string]interface{}{
		"subId": subId,
	})
}

// getAllSubResponse defines the response structure for listing all subscriptions
// Wraps the subscription array in a "data" field for consistent JSON response format
type getAllSubResponse struct {
	Data []*models.Subscription `json:"data"` // Array of subscription objects
}

// getAllSubscriptions handles HTTP GET request to retrieve all subscriptions
// This endpoint returns a list of all subscriptions in the system
func (h *Handler) getAllSubscriptions(c *gin.Context) {
	// Retrieve all subscriptions from the service layer
	subs, err := h.services.SubscriptionStore.GetAll()
	if err != nil {
		// Return 500 Internal Server Error if data retrieval fails
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Return 200 OK with subscriptions wrapped in a standardized response structure
	c.JSON(http.StatusOK, getAllSubResponse{
		Data: subs,
	})
}

func (h *Handler) getSubscriptionById(c *gin.Context) {

}

func (h *Handler) updateSubscription(c *gin.Context) {

}

func (h *Handler) deleteSubscription(c *gin.Context) {

}
