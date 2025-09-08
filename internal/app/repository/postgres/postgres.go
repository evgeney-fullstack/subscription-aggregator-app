package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	subscriptionTable = "subscriptions" // Database table name for subscriptions
)

// Config holds PostgreSQL connection configuration parameters
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL database connection with context support
// Returns sqlx.DB instance or error if connection fails
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	var db *sqlx.DB
	var err error

	// Retry logic for connection
	for i := 0; i < 3; i++ {
		db, err = sqlx.Open("postgres", connectionString)
		if err != nil {
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		if err = db.PingContext(ctx); err != nil {
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		break
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect after retries: %w", err)
	}

	return db, nil
}
