package handler

import (
	"net/http"
	"strconv"

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

// getSubscriptionById handles HTTP GET request to retrieve a specific subscription by ID
// This endpoint returns a single subscription based on the provided subscription_id parameter
func (h *Handler) getSubscriptionById(c *gin.Context) {
	// Extract and convert subscription_id parameter from URL path to integer
	// The parameter is expected to be in the format: /subscriptions/{subscription_id}
	subID, err := strconv.Atoi(c.Param("subscription_id"))
	if err != nil {
		// Return 400 Bad Request if the parameter is not a valid integer
		newErrorResponse(c, http.StatusBadRequest, "invalid subscription_id param")
		return
	}

	// Retrieve the subscription from the service layer using the extracted ID
	sub, err := h.services.SubscriptionStore.GetById(subID)
	if err != nil {
		// Return 500 Internal Server Error if data retrieval fails
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Return 200 OK with the subscription data in JSON format
	c.JSON(http.StatusOK, sub)
}

func (h *Handler) updateSubscription(c *gin.Context) {

	// Extract and convert subscription_id parameter from URL path to integer
	// The parameter is expected to be in the format: /subscriptions/{subscription_id}
	subID, err := strconv.Atoi(c.Param("subscription_id"))
	if err != nil {
		// Return 400 Bad Request if the parameter is not a valid integer
		newErrorResponse(c, http.StatusBadRequest, "invalid subscription_id param")
		return
	}

	var input models.UpdateSubscription
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.SubscriptionStore.Update(subID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Operation completed successfully",
	})
}

func (h *Handler) deleteSubscription(c *gin.Context) {
	// Extract and convert subscription_id parameter from URL path to integer
	// The parameter is expected to be in the format: /subscriptions/{subscription_id}
	subID, err := strconv.Atoi(c.Param("subscription_id"))
	if err != nil {
		// Return 400 Bad Request if the parameter is not a valid integer
		newErrorResponse(c, http.StatusBadRequest, "invalid subscription_id param")
		return
	}

	// Service will delete the subscription
	err = h.services.SubscriptionStore.Delete(subID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Return HTTP 200 OK status with a success confirmation message
	// This response is typically used for operations that don't require returning data,
	// but need to confirm successful execution to the client
	c.JSON(http.StatusOK, statusResponse{
		Status: "Operation completed successfully",
	})

}

func (h *Handler) getSubscriptionSummary(c *gin.Context) {

	var filter models.SubscriptionFilter
	if err := c.BindJSON(&filter); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	totalCost, err := h.services.SubscriptionStore.GetSubscriptionSummary(filter)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	filter.TotalCost = totalCost
	filter.Currency = "RUB"

	c.JSON(http.StatusOK, filter)
}
