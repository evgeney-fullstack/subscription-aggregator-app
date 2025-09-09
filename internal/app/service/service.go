package service

import (
	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/models"
	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/repository/postgres"
)

// SubscriptionStore defines business logic operations for subscriptions
type SubscriptionStore interface {
	Create(sub models.Subscription) (int, error)
	GetAll() ([]*models.Subscription, error)
	GetById(subID int) (models.Subscription, error)
	Delete(subID int) error
	Update(subID int, input models.UpdateSubscription) error
	GetSubscriptionSummary(filter models.SubscriptionFilter) (int, error)
}

// Service layer aggregates all business logic services
type Service struct {
	SubscriptionStore
}

// NewService constructs new Service layer with business logic
func NewService(repos *postgres.Repository) *Service {
	return &Service{
		SubscriptionStore: NewSubscriptionService(repos.SubscriptionStore),
	}
}
