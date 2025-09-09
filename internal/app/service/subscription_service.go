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

	// Parse string UserID from API request into UUID format for database storage
	_, err := uuid.Parse(sub.UserID)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Parse string date from API request into time.Time object
	// Uses "01-2006" format (month-year) following Go's reference date format
	_, err = time.Parse("01-2006", sub.StartDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date format, expected MM-YYYY: %w", err)
	}

	// Delegate to repository layer for actual database persistence
	return s.repo.Create(sub)
}

// ConvertDBToAPIModel transforms a database model to an API response model
func сonvertDBToAPIModel(subdb models.SubscriptionDB) models.Subscription {
	return models.Subscription{
		Id:          subdb.Id,
		ServiceName: subdb.ServiceName,
		Price:       subdb.Price,
		UserID:      subdb.UserID.String(),
		StartDate:   subdb.StartDate.Format("01-2006"),
		FinishDate:  subdb.FinishDate.Format("01-2006"),
	}
}

// GetAll retrieves all subscriptions from the repository and converts them to API model format
// Returns a slice of Subscription models or an error if data retrieval fails
func (s *SubscriptionService) GetAll() ([]*models.Subscription, error) {
	// Retrieve all subscriptions from the repository layer (database)
	subsDB, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve subscriptions from repository: %w", err)
	}

	// Initialize slice for API response models
	// Pre-allocate capacity for better performance with large datasets
	subs := make([]*models.Subscription, 0, len(subsDB))

	// Convert each database model to API response model
	for i := range subsDB {
		sub := сonvertDBToAPIModel(subsDB[i])
		subs = append(subs, &sub)
	}

	return subs, nil
}

// GetById implements business logic for retrieving subscription by ID (to be implemented)
func (s *SubscriptionService) GetById(subID int) (models.Subscription, error) {
	var sub models.Subscription

	// Retrieve  subscription by ID from the repository layer (database)
	subDB, err := s.repo.GetById(subID)
	if err != nil {
		return sub, fmt.Errorf("failed to retrieve subscriptions from repository: %w", err)
	}

	// Convert each database model to API response model
	sub = сonvertDBToAPIModel(subDB)

	return sub, nil
}

// Delete implements subscription deletion business logic (to be implemented)
func (s *SubscriptionService) Delete(subID int) error {
	return s.repo.Delete(subID)
}

// Update handles the business logic for updating an existing subscription
// Validates input data before delegating to the repository layer for persistence
func (s *SubscriptionService) Update(subID int, input models.UpdateSubscription) error {
	// Validate input data using the model's validation method
	// This ensures business rules are enforced before database operations
	if err := input.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if input.StartDate != nil {
		// Parse string date from API request into time.Time object
		// Uses "01-2006" format (month-year) following Go's reference date format
		_, err := time.Parse("01-2006", *input.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start date format, expected MM-YYYY: %w", err)
		}
	}

	// Delegate the update operation to the repository layer
	// The repository handles the actual database interaction
	return s.repo.Update(subID, input)
}

// GetSubscriptionSummary converts API filters to DB format and calculates total cost
func (s *SubscriptionService) GetSubscriptionSummary(filter models.SubscriptionFilter) (int, error) {
	return s.repo.GetSubscriptionSummary(filter)
}
