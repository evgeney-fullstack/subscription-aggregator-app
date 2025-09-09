package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Subscription represents the subscription model for API requests/responses
// Used for JSON marshaling/unmarshaling with string-based date fields
type Subscription struct {
	Id          int    `json:"id" db:"id"`                      // Unique identifier
	ServiceName string `json:"service_name" binding:"required"` // Name of the service (required)
	Price       int    `json:"price" binding:"required"`        // Subscription price (required)
	UserID      string `json:"user_id" binding:"required"`      // User identifier as string (required)
	StartDate   string `json:"start_date" binding:"required"`   // Start date in string format (required)
	FinishDate  string `json:"finish_date"`                     // End date as timestamp
}

// SubscriptionDB represents the subscription model for database operations
// Uses proper data types for database storage (UUID, time.Time)
type SubscriptionDB struct {
	Id          int       `db:"id"`           // Unique identifier
	ServiceName string    `db:"service_name"` // Name of the service
	Price       int       `db:"price"`        // Subscription price
	UserID      uuid.UUID `db:"user_id"`      // User identifier as UUID
	StartDate   time.Time `db:"start_date"`   // Start date as timestamp
	FinishDate  time.Time `db:"finish_date"`  // End date as timestamp
}

// UpdateSubscription defines the structure for subscription update requests
// Uses pointer fields to distinguish between missing values and zero values
// This allows for partial updates (PATCH semantics) where only provided fields are updated
type UpdateSubscription struct {
	Price     *int    `json:"price" `     // Optional new price value (pointer allows nil for no update)
	StartDate *string `json:"start_date"` // Optional new start date in "MM-YYYY" format
}

// Validate ensures the update request contains at least one field to update
// Prevents empty update operations that would make no changes to the resource
func (i UpdateSubscription) Validate() error {
	// Check that at least one field is provided for update
	if i.Price == nil && i.StartDate == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

// SubscriptionFilter defines the request/response structure for subscription summary API
// Used for both input parameters and output response
type SubscriptionFilter struct {
	TotalCost int     `json:"total_cost"` // Calculated total cost of subscriptions (output only)
	Currency  string  `json:"currency"`   // Currency code for the total cost (output only)
	Period    Period  `json:"period"`     // Time period for the summary calculation (input)
	Filters   Filters `json:"filters"`    // Optional filters for the summary (input)
}

// Period defines the time range for subscription summary calculation
type Period struct {
	StartDate  string `json:"start_date" binding:"required"`  // Start date in MM-YYYY format
	FinishDate string `json:"finish_date" binding:"required"` // End date in MM-YYYY format
}

// Filters contains optional criteria to narrow down subscription summary
type Filters struct {
	UserID      *string `json:"user_id"`      // Optional filter by user ID
	ServiceName *string `json:"service_name"` // Optional filter by service name
}

// SubscriptionFilterDB is the database representation of subscription filters
// Used to pass filter criteria to repository layer
type SubscriptionFilterDB struct {
	ServiceName *string   `db:"service_name"` // Service name filter (optional)
	UserID      *string   `db:"user_id"`      // User ID filter (optional)
	StartDate   time.Time `db:"start_date"`   // Start date in time.Time format
	FinishDate  time.Time `db:"finish_date"`  // End date in time.Time format
}
