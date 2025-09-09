package postgres

import (
	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/models"
	"github.com/jmoiron/sqlx"
)

// SubscriptionStore defines CRUD operations for subscription management
type SubscriptionStore interface {
	Create(sub models.SubscriptionDB) (int, error)
	GetAll() ([]models.SubscriptionDB, error)
	GetById(subID int) (models.SubscriptionDB, error)
	Delete()
	Update()
}

// Repository aggregates all store interfaces for database operations
type Repository struct {
	SubscriptionStore
}

// NewRepository constructs a new Repository with all available stores
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		SubscriptionStore: NewSubscriptionRepository(db),
	}
}
