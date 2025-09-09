package postgres

import (
	"errors"
	"fmt"
	"strings"
	"time"

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
func (r *SubscriptionRepository) GetAll() ([]models.SubscriptionDB, error) {
	var subDB []models.SubscriptionDB

	query := fmt.Sprintf("SELECT * FROM %s", subscriptionTable)
	err := r.db.Select(&subDB, query)

	return subDB, err
}

// GetById implements retrieval of subscription by ID (to be implemented)
func (r *SubscriptionRepository) GetById(subID int) (models.SubscriptionDB, error) {

	var subDB models.SubscriptionDB

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", subscriptionTable)
	err := r.db.Get(&subDB, query, subID)

	return subDB, err
}

// Delete implements subscription deletion logic (to be implemented)
func (r *SubscriptionRepository) Delete(subID int) error {

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", subscriptionTable)
	result, err := r.db.Exec(query, subID)
	if err != nil {
		return err
	}

	//Check if the card has been deleted
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("card not found")
	}
	return nil

}

// Update implements subscription update logic with partial update support
// Handles dynamic SQL query generation based on provided fields
func (r *SubscriptionRepository) Update(subID int, input models.UpdateSubscription) error {
	// Initialize slices for building dynamic SET clause and arguments
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1 // Positional parameter counter

	// Handle price update if provided
	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *input.Price)
		argId++
	}

	// Handle start date update if provided
	if input.StartDate != nil {
		setValues = append(setValues, fmt.Sprintf("start_date=$%d", argId))

		// Parse string date to time.Time for database storage
		startData, err := time.Parse("01-2006", *input.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start date format, expected MM-YYYY: %w", err)
		}
		args = append(args, startData)
		argId++

		// Calculate subscription end date (1 month duration from start date)
		finishDate := startData.AddDate(0, 1, 0)
		setValues = append(setValues, fmt.Sprintf("finish_date=$%d", argId))

		args = append(args, finishDate)
		argId++
	}

	// Join SET clauses with commas
	setQuery := strings.Join(setValues, ", ")

	// Build final SQL query with WHERE clause
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", subscriptionTable, setQuery, argId)

	// Add subscription ID as the last parameter
	args = append(args, subID)

	// Execute the query
	_, err := r.db.Exec(query, args...)
	return err
}
