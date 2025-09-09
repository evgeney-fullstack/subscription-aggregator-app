package models

import (
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
	FinishDate  string `json:"finish_date"`                     // End date in string format (optional)
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
