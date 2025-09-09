package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/handler"
	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/repository/postgres"
	"github.com/evgeney-fullstack/subscription-aggregator-app/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestMain устанавливает тестовое окружение
func TestMain(m *testing.M) {
	// Устанавливаем режим Gin в тестовый
	gin.SetMode(gin.TestMode)

	// Запускаем тесты
	m.Run()
}

// setupTestContainer настраивает контейнер PostgreSQL для тестирования
func setupTestContainer(ctx context.Context) (postgres.Config, func(), error) {
	// Запуск контейнера PostgreSQL с использованием GenericContainer
	postgresReq := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(30 * time.Second),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	if err != nil {
		return postgres.Config{}, nil, fmt.Errorf("failed to start PostgreSQL container: %w", err)
	}

	// Получение параметров подключения к PostgreSQL
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return postgres.Config{}, nil, err
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return postgres.Config{}, nil, err
	}

	// Конфигурация для подключения к тестовому контейнеру PostgreSQL
	postgresCfg := postgres.Config{
		Host:     host,
		Port:     port.Port(),
		Username: "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	// Функция очистки для остановки контейнера
	cleanup := func() {
		if terminateErr := postgresContainer.Terminate(ctx); terminateErr != nil {
			log.Printf("failed to terminate PostgreSQL container: %v", terminateErr)
		}
	}

	return postgresCfg, cleanup, nil
}

// setupTestServer создает и настраивает тестовый сервер
func setupTestServer(postgresCfg postgres.Config) (*gin.Engine, error) {
	// Инициализация PostgreSQL
	db, err := postgres.NewPostgresDB(postgresCfg)
	if err != nil {
		return nil, err
	}

	// Добавьте создание таблицы
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS subscriptions (
            id SERIAL PRIMARY KEY,
            service_name VARCHAR(255) NOT NULL,
            price INTEGER NOT NULL,
            user_id UUID NOT NULL,
            start_date TIMESTAMP NOT NULL,
            finish_date TIMESTAMP NOT NULL
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Инициализация репозиториев
	repos := postgres.NewRepository(db)

	// Инициализация сервисов
	services := service.NewService(repos)

	// Инициализация обработчиков
	handler := handler.NewHandler(services)

	// Настройка маршрутов
	router := handler.InitRoutes()

	return router, nil
}

// TestSignUpIntegration is testing the endpoint of user registration
func TestСreateSubscriptionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	ctx := context.Background()

	dbConfig, cleanup, err := setupTestContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to set up test container: %v", err)
	}
	defer cleanup()

	router, err := setupTestServer(dbConfig)
	if err != nil {
		t.Fatalf("Failed to set up test server: %v", err)
	}

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful creation",
			payload: map[string]interface{}{
				"service_name": "TEST",
				"price":        123,
				"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				"start_date":   "07-2025",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "subId")
			},
		},
		{
			name: "Error -  invalid user ID format",
			payload: map[string]interface{}{
				"service_name": "TEST",
				"price":        123,
				"user_id":      "60601fee",
				"start_date":   "07-2025",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "invalid user ID format")
			},
		},
		{
			name: "Error -  invalid start date format",
			payload: map[string]interface{}{
				"service_name": "TEST",
				"price":        123,
				"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				"start_date":   "01-07-2025",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "invalid start date format")
			},
		},
		{
			name: "Error -  service_name empty",
			payload: map[string]interface{}{
				"service_name": "",
				"price":        123,
				"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				"start_date":   "07-2025",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "Error")
			},
		},
		{
			name: "Error -  price empty",
			payload: map[string]interface{}{
				"service_name": "TEST",
				"price":        0,
				"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				"start_date":   "07-2025",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "Error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request preparation
			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/subscriptions/", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Request execution
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			// Checking the response status
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			// Checking the response body
			if tt.checkResponse != nil {
				tt.checkResponse(t, recorder)
			}
		})
	}
}

func TestGetAllSubscriptionsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	ctx := context.Background()

	dbConfig, cleanup, err := setupTestContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to set up test container: %v", err)
	}
	defer cleanup()

	router, err := setupTestServer(dbConfig)
	if err != nil {
		t.Fatalf("Failed to set up test server: %v", err)
	}

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful creation TEST1",
			payload: map[string]interface{}{
				"service_name": "TEST1",
				"price":        200,
				"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				"start_date":   "06-2025",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "subId")
			},
		},
		{
			name: "Successful creation TEST2",
			payload: map[string]interface{}{
				"service_name": "TEST2",
				"price":        300,
				"user_id":      "60691fee-2bf1-4721-ae6f-7036e79a0cba",
				"start_date":   "07-2025",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "subId")
			},
		},
		{
			name:           "Successful getAllSubscriptions",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "data")
			},
		},
	}

	cunter := 0

	for _, tt := range tests {

		if cunter < 2 {
			t.Run(tt.name, func(t *testing.T) {
				// Request preparation
				body, _ := json.Marshal(tt.payload)
				req, _ := http.NewRequest("POST", "/subscriptions/", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				// Request execution
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, req)

				// Checking the response status
				assert.Equal(t, tt.expectedStatus, recorder.Code)

				// Checking the response body
				if tt.checkResponse != nil {
					tt.checkResponse(t, recorder)
				}
			})
			cunter++
		} else {
			t.Run(tt.name, func(t *testing.T) {
				// Request preparation
				body, _ := json.Marshal(tt.payload)
				req, _ := http.NewRequest("GET", "/subscriptions/", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				// Request execution
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, req)

				// Checking the response status
				assert.Equal(t, tt.expectedStatus, recorder.Code)

				// Checking the response body
				if tt.checkResponse != nil {
					tt.checkResponse(t, recorder)
				}
			})

		}
	}
}

func TestGetSubscriptionByIdIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	ctx := context.Background()

	dbConfig, cleanup, err := setupTestContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to set up test container: %v", err)
	}
	defer cleanup()

	router, err := setupTestServer(dbConfig)
	if err != nil {
		t.Fatalf("Failed to set up test server: %v", err)
	}

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful creation TEST1",
			payload: map[string]interface{}{
				"service_name": "TEST1",
				"price":        200,
				"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				"start_date":   "06-2025",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "subId")
			},
		},
		{
			name: "Successful creation TEST2",
			payload: map[string]interface{}{
				"service_name": "TEST2",
				"price":        300,
				"user_id":      "60691fee-2bf1-4721-ae6f-7036e79a0cba",
				"start_date":   "07-2025",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "subId")
			},
		},
		{
			name:           "Successful getAllSubscriptions",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "TEST2")
			},
		},
	}

	cunter := 0

	for _, tt := range tests {

		if cunter < 2 {
			t.Run(tt.name, func(t *testing.T) {
				// Request preparation
				body, _ := json.Marshal(tt.payload)
				req, _ := http.NewRequest("POST", "/subscriptions/", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				// Request execution
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, req)

				// Checking the response status
				assert.Equal(t, tt.expectedStatus, recorder.Code)

				// Checking the response body
				if tt.checkResponse != nil {
					tt.checkResponse(t, recorder)
				}
			})
			cunter++
		} else {
			t.Run(tt.name, func(t *testing.T) {
				// Request preparation
				body, _ := json.Marshal(tt.payload)
				req, _ := http.NewRequest("GET", "/subscriptions/2", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				// Request execution
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, req)

				// Checking the response status
				assert.Equal(t, tt.expectedStatus, recorder.Code)

				// Checking the response body
				if tt.checkResponse != nil {
					tt.checkResponse(t, recorder)
				}
			})

		}
	}
}

func TestDeleteSubscriptionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	ctx := context.Background()

	dbConfig, cleanup, err := setupTestContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to set up test container: %v", err)
	}
	defer cleanup()

	router, err := setupTestServer(dbConfig)
	if err != nil {
		t.Fatalf("Failed to set up test server: %v", err)
	}

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful creation TEST1",
			payload: map[string]interface{}{
				"service_name": "TEST1",
				"price":        200,
				"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				"start_date":   "06-2025",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "subId")
			},
		},
		{
			name: "Successful creation TEST2",
			payload: map[string]interface{}{
				"service_name": "TEST2",
				"price":        300,
				"user_id":      "60691fee-2bf1-4721-ae6f-7036e79a0cba",
				"start_date":   "07-2025",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "subId")
			},
		},
		{
			name:           "Successful delete subscription",
			payload:        map[string]interface{}{},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "Operation completed successfully")
			},
		},
	}

	cunter := 0

	for _, tt := range tests {

		if cunter < 2 {
			t.Run(tt.name, func(t *testing.T) {
				// Request preparation
				body, _ := json.Marshal(tt.payload)
				req, _ := http.NewRequest("POST", "/subscriptions/", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				// Request execution
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, req)

				// Checking the response status
				assert.Equal(t, tt.expectedStatus, recorder.Code)

				// Checking the response body
				if tt.checkResponse != nil {
					tt.checkResponse(t, recorder)
				}
			})
			cunter++
		} else {
			t.Run(tt.name, func(t *testing.T) {
				// Request preparation
				body, _ := json.Marshal(tt.payload)
				req, _ := http.NewRequest("DELETE", "/subscriptions/2", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				// Request execution
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, req)

				// Checking the response status
				assert.Equal(t, tt.expectedStatus, recorder.Code)

				// Checking the response body
				if tt.checkResponse != nil {
					tt.checkResponse(t, recorder)
				}
			})

		}
	}
}

func TestUpdateSubscriptionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	ctx := context.Background()

	dbConfig, cleanup, err := setupTestContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to set up test container: %v", err)
	}
	defer cleanup()

	router, err := setupTestServer(dbConfig)
	if err != nil {
		t.Fatalf("Failed to set up test server: %v", err)
	}

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Successful creation TEST1",
			payload: map[string]interface{}{
				"service_name": "TEST1",
				"price":        200,
				"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
				"start_date":   "06-2025",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "subId")
			},
		},
		{
			name: "Successful update subscription",
			payload: map[string]interface{}{
				"price":      100,
				"start_date": "04-2025",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "Operation completed successfully")
			},
		},
	}

	cunter := 0

	for _, tt := range tests {

		if cunter < 1 {
			t.Run(tt.name, func(t *testing.T) {
				// Request preparation
				body, _ := json.Marshal(tt.payload)
				req, _ := http.NewRequest("POST", "/subscriptions/", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				// Request execution
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, req)

				// Checking the response status
				assert.Equal(t, tt.expectedStatus, recorder.Code)

				// Checking the response body
				if tt.checkResponse != nil {
					tt.checkResponse(t, recorder)
				}
			})
			cunter++
		} else {
			t.Run(tt.name, func(t *testing.T) {

				// Request preparation
				body, _ := json.Marshal(tt.payload)
				req, _ := http.NewRequest("PUT", "/subscriptions/1", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				// Request execution
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, req)

				// Checking the response status
				assert.Equal(t, tt.expectedStatus, recorder.Code)

				// Checking the response body
				if tt.checkResponse != nil {
					tt.checkResponse(t, recorder)
				}
			})

		}
	}
}
