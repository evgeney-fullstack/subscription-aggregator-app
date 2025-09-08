package postgres

import (
	"fmt"

	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/models"
	"github.com/jmoiron/sqlx"
)

// SubscriptionRepository implements SubscriptionStore for PostgreSQL
type SubscriptionRepository struct {
	db *sqlx.DB
}

// NewSubscriptionRepository creates a new subscription repository instance
func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create inserts a new subscription record into the database
// Returns the ID of the newly created subscription or an error
func (r *SubscriptionRepository) Create(subDB models.SubscriptionDB) (int, error) {
	// Begin a database transaction to ensure atomic operation
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var subId int
	// Prepare SQL query for subscription insertion with parameter binding
	// Uses RETURNING clause to get the auto-generated ID
	createSubQuery := fmt.Sprintf("INSERT INTO %s (service_name, price, user_id, start_date, finish_date) VALUES ($1, $2, $3, $4, $5) RETURNING id", subscriptionTable)

	// Execute the query within the transaction
	row := tx.QueryRow(createSubQuery, subDB.ServiceName, subDB.Price, subDB.UserID, subDB.StartDate, subDB.FinishDate)

	// Retrieve the auto-generated ID from the result
	if err := row.Scan(&subId); err != nil {
		// Rollback transaction in case of error to maintain data consistency
		tx.Rollback()
		return 0, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Commit the transaction to persist changes
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return subId, nil
}

// GetAll implements retrieval of all subscriptions (to be implemented)
func (r *SubscriptionRepository) GetAll() {

}

// GetById implements retrieval of subscription by ID (to be implemented)
func (r *SubscriptionRepository) GetById() {

}

// Delete implements subscription deletion logic (to be implemented)
func (r *SubscriptionRepository) Delete() {

}

// Update implements subscription update logic (to be implemented)
func (r *SubscriptionRepository) Update() {

}
