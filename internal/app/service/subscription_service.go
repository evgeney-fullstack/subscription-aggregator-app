package service

import (
	"fmt"
	"time"

	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/models"
	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/repository/postgres"
	"github.com/google/uuid"
)

// SubscriptionService implements business logic for subscription operations
type SubscriptionService struct {
	repo postgres.SubscriptionStore
}

// NewSubscriptionService creates a new subscription service instance
func NewSubscriptionService(repo postgres.SubscriptionStore) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

// Create handles the business logic for creating a new subscription
// Transforms API model (Subscription) to database model (SubscriptionDB)
// Performs data validation and transformation before persistence
func (s *SubscriptionService) Create(sub models.Subscription) (int, error) {
	var subDB models.SubscriptionDB

	// Parse string UserID from API request into UUID format for database storage
	userID, err := uuid.Parse(sub.UserID)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Parse string date from API request into time.Time object
	// Uses "01-2006" format (month-year) following Go's reference date format
	startData, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date format, expected MM-YYYY: %w", err)
	}

	// Map fields from API model to database model
	subDB.Price = sub.Price
	subDB.ServiceName = sub.ServiceName
	subDB.UserID = userID
	subDB.StartDate = startData

	// Calculate subscription end date (1 month duration from start date)
	subDB.FinishDate = startData.AddDate(0, 1, 0)

	// Delegate to repository layer for actual database persistence
	return s.repo.Create(subDB)
}

// GetAll implements business logic for retrieving all subscriptions (to be implemented)
func (s *SubscriptionService) GetAll() {

}

// GetById implements business logic for retrieving subscription by ID (to be implemented)
func (s *SubscriptionService) GetById() {

}

// Delete implements subscription deletion business logic (to be implemented)
func (s *SubscriptionService) Delete() {

}

// Update implements subscription update business logic (to be implemented)
func (s *SubscriptionService) Update() {

}
