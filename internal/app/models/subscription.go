package models

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a user's subscription to a service in the system
// Contains information about subscription details, pricing and duration
type Subscription struct {
	Id          int       `json:"-" db:"id"`                       // Unique identifier (not exposed in JSON)
	ServiceName string    `json:"service_name" binding:"required"` // Name of the subscribed service (required)
	Price       int       `json:"price" binding:"required"`        // Subscription price in currency units (required)
	UserID      uuid.UUID `json:"user_id" binding:"required"`      // Unique identifier of the user who owns the subscription (required)
	StartDate   time.Time `json:"start_date" binding:"required"`   // Date when the subscription becomes active (required)
	FinishDate  time.Time `json:"finish_date"`                     // Date when the subscription expires (optional)
}
