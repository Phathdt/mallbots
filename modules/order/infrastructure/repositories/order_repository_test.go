package repositories

import (
	"context"
	"fmt"
	"mallbots/modules/order/domain/constants"
	"mallbots/modules/order/domain/entities"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phathdt/service-context/core"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func createContainer(t *testing.T) (*postgres.PostgresContainer, error) {
	ctx := context.Background()
	dbUsername := "postgres"
	dbPassword := "123123123"
	dbName := "mallbots_test"

	schemaFile := filepath.Join("../../../../schema.gen.sql")
	seedFile := filepath.Join("../../../../seed.sql")

	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithInitScripts(schemaFile, seedFile),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUsername),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %v", err)
		}
	})

	return postgresContainer, nil
}

func createTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	container, err := createContainer(t)
	require.NoError(t, err, "failed to create container")

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err, "failed to get connection string")

	poolConfig, err := pgxpool.ParseConfig(connStr)
	require.NoError(t, err, "failed to parse connection string")

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	require.NoError(t, err, "failed to create connection pool")

	err = pool.Ping(ctx)
	require.NoError(t, err, "failed to ping database")

	return pool
}

// createTestUsers creates test users in the database
func createTestUsers(ctx context.Context, db *pgxpool.Pool) error {
	testUsers := []struct {
		email    string
		password string
		fullName string
	}{
		{"test1@example.com", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", "Test User 1"},
		{"test2@example.com", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", "Test User 2"},
		{"test3@example.com", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", "Test User 3"},
		{"test4@example.com", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", "Test User 4"},
	}

	for _, user := range testUsers {
		_, err := db.Exec(ctx,
			"INSERT INTO users (email, password, full_name, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())",
			user.email, user.password, user.fullName)
		if err != nil {
			return fmt.Errorf("failed to create test user: %w", err)
		}
	}

	return nil
}

func TestOrderRepository(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	err := createTestUsers(ctx, db)
	require.NoError(t, err, "failed to create test users")

	repo := NewOrderRepository(db)

	t.Run("Create Order with Items", func(t *testing.T) {
		// Create order
		order := &entities.Order{
			UserID:          1,
			Status:          constants.OrderStatusPending,
			PaymentStatus:   constants.PaymentStatusPending,
			TotalAmount:     100.00,
			ShippingAddress: "123 Test St",
			ShippingCity:    "Test City",
			ShippingCountry: "Test Country",
			ShippingZip:     "12345",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		// Create order
		createdOrder, err := repo.Create(ctx, order)
		require.NoError(t, err)
		require.NotZero(t, createdOrder.ID)
		require.Equal(t, order.UserID, createdOrder.UserID)
		require.Equal(t, order.TotalAmount, createdOrder.TotalAmount)
		require.Equal(t, order.Status, createdOrder.Status)
		require.Equal(t, order.PaymentStatus, createdOrder.PaymentStatus)

		// Create order items
		items := []*entities.OrderItem{
			{
				OrderID:   createdOrder.ID,
				ProductID: 1,
				Quantity:  2,
				Price:     25.00,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				OrderID:   createdOrder.ID,
				ProductID: 2,
				Quantity:  1,
				Price:     50.00,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		err = repo.CreateOrderItems(ctx, createdOrder.ID, items)
		require.NoError(t, err)

		// Get order with items
		fetchedOrder, err := repo.GetByID(ctx, createdOrder.ID)
		require.NoError(t, err)
		require.Equal(t, createdOrder.ID, fetchedOrder.ID)
		require.Len(t, fetchedOrder.Items, 2)
		require.Equal(t, items[0].ProductID, fetchedOrder.Items[0].ProductID)
		require.Equal(t, items[0].Quantity, fetchedOrder.Items[0].Quantity)
		require.Equal(t, items[0].Price, fetchedOrder.Items[0].Price)
	})

	t.Run("Get User Orders with Pagination", func(t *testing.T) {
		userID := int32(2)

		// Create multiple orders for user
		orders := []*entities.Order{
			{
				UserID:          userID,
				Status:          constants.OrderStatusPending,
				PaymentStatus:   constants.PaymentStatusPending,
				TotalAmount:     100.00,
				ShippingAddress: "123 Test St",
				ShippingCity:    "Test City",
				ShippingCountry: "Test Country",
				ShippingZip:     "12345",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				UserID:          userID,
				Status:          constants.OrderStatusConfirmed,
				PaymentStatus:   constants.PaymentStatusPaid,
				TotalAmount:     200.00,
				ShippingAddress: "456 Test St",
				ShippingCity:    "Test City",
				ShippingCountry: "Test Country",
				ShippingZip:     "12345",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		}

		for _, order := range orders {
			createdOrder, err := repo.Create(ctx, order)
			require.NoError(t, err)

			// Add items to each order
			items := []*entities.OrderItem{
				{
					OrderID:   createdOrder.ID,
					ProductID: 1,
					Quantity:  1,
					Price:     50.00,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}
			err = repo.CreateOrderItems(ctx, createdOrder.ID, items)
			require.NoError(t, err)
		}

		// Test pagination
		paging := &core.Paging{
			Page:  1,
			Limit: 1,
		}

		// Get first page
		userOrders, err := repo.GetByUserID(ctx, userID, paging)
		require.NoError(t, err)
		require.Len(t, userOrders, 1)
		require.Equal(t, int64(2), paging.Total)

		// Get second page
		paging.Page = 2
		userOrders, err = repo.GetByUserID(ctx, userID, paging)
		require.NoError(t, err)
		require.Len(t, userOrders, 1)
	})

	t.Run("Update Order Status", func(t *testing.T) {
		// Create initial order
		order := &entities.Order{
			UserID:          3,
			Status:          constants.OrderStatusPending,
			PaymentStatus:   constants.PaymentStatusPending,
			TotalAmount:     100.00,
			ShippingAddress: "123 Test St",
			ShippingCity:    "Test City",
			ShippingCountry: "Test Country",
			ShippingZip:     "12345",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		createdOrder, err := repo.Create(ctx, order)
		require.NoError(t, err)

		// Update order status
		err = repo.UpdateStatus(ctx, createdOrder.ID, constants.OrderStatusConfirmed)
		require.NoError(t, err)

		// Verify update
		updatedOrder, err := repo.GetByID(ctx, createdOrder.ID)
		require.NoError(t, err)
		require.Equal(t, constants.OrderStatusConfirmed, updatedOrder.Status)
	})

	t.Run("Update Payment Status", func(t *testing.T) {
		// Create initial order
		order := &entities.Order{
			UserID:          4,
			Status:          constants.OrderStatusConfirmed,
			PaymentStatus:   constants.PaymentStatusPending,
			TotalAmount:     100.00,
			ShippingAddress: "123 Test St",
			ShippingCity:    "Test City",
			ShippingCountry: "Test Country",
			ShippingZip:     "12345",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		createdOrder, err := repo.Create(ctx, order)
		require.NoError(t, err)

		// Update payment status
		err = repo.UpdatePaymentStatus(ctx, createdOrder.ID, constants.PaymentStatusPaid)
		require.NoError(t, err)

		// Verify update
		updatedOrder, err := repo.GetByID(ctx, createdOrder.ID)
		require.NoError(t, err)
		require.Equal(t, constants.PaymentStatusPaid, updatedOrder.PaymentStatus)
	})

	t.Run("Get Non-existent Order", func(t *testing.T) {
		_, err := repo.GetByID(ctx, 99999)
		require.Error(t, err)
	})
}
