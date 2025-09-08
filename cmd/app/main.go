package main

import (
	"os"

	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/handler"
	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/repository/postgres"
	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/server"
	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	// Configuring the logs format in JSON for better structuring and compatibility
	// with monitoring systems (Kibana, Elasticsearch, etc.)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Loading environment variables from the config.env file
	if err := godotenv.Load("config.env"); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// Initializing a connection to PostgreSQL using parameters from environment variables
	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// Initializing repositories for working with data
	// repos provides access to PostgreSQL data
	repos := postgres.NewRepository(db)

	// Creating a service layer with dependency injection
	// service encapsulates the business logic of the application
	service := service.NewService(repos)

	// Initialization of HTTP handlers with the introduction of a service layer
	// Handlers will use business logic via service
	handlers := handler.NewHandler(service)

	// Creating a server instance
	srv := new(server.Server)

	// Launching an HTTPS server with configuration from environment variables
	// Using HOST and HOST_PORT from config.env
	if err := srv.Run(os.Getenv("HOST"), os.Getenv("HOST_PORT"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	}

}
