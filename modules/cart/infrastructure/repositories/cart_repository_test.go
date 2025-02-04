package repositories

import (
	"context"
	"fmt"
	"mallbots/modules/cart/domain/entities"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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
	// Create test users
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

func TestCartRepository(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	ctx := context.Background()
	err := createTestUsers(ctx, db)
	require.NoError(t, err, "failed to create test users")

	repo := NewCartRepository(db)

	t.Run("Create and Get Cart Item", func(t *testing.T) {
		// Create test cart item
		item := &entities.CartItem{
			UserID:    1,
			ProductID: 1,
			Quantity:  2,
			Price:     10.99,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create item
		createdItem, err := repo.Create(ctx, item)
		require.NoError(t, err)
		require.NotZero(t, createdItem.ID)
		require.Equal(t, item.UserID, createdItem.UserID)
		require.Equal(t, item.ProductID, createdItem.ProductID)

		// Get item
		fetchedItem, err := repo.GetByUserAndProduct(ctx, item.UserID, item.ProductID)
		require.NoError(t, err)
		require.Equal(t, createdItem.ID, fetchedItem.ID)
	})

	t.Run("Update Cart Item", func(t *testing.T) {
		// Create initial item
		item := &entities.CartItem{
			UserID:    2,
			ProductID: 1,
			Quantity:  1,
			Price:     10.99,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		createdItem, err := repo.Create(ctx, item)
		require.NoError(t, err)

		// Update quantity
		createdItem.Quantity = 3
		createdItem.UpdatedAt = time.Now()

		err = repo.Update(ctx, createdItem)
		require.NoError(t, err)

		// Verify update
		updatedItem, err := repo.GetByUserAndProduct(ctx, createdItem.UserID, createdItem.ProductID)
		require.NoError(t, err)
		require.Equal(t, int32(3), updatedItem.Quantity)
	})

	t.Run("Delete Cart Item", func(t *testing.T) {
		// Create item to delete
		item := &entities.CartItem{
			UserID:    3,
			ProductID: 1,
			Quantity:  1,
			Price:     10.99,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		createdItem, err := repo.Create(ctx, item)
		require.NoError(t, err)

		// Delete item
		err = repo.Delete(ctx, createdItem.UserID, createdItem.ProductID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByUserAndProduct(ctx, createdItem.UserID, createdItem.ProductID)
		require.Error(t, err)
	})

	t.Run("Delete All User Cart Items", func(t *testing.T) {
		userID := int32(4)

		// Create multiple items for the user
		items := []*entities.CartItem{
			{
				UserID:    userID,
				ProductID: 1,
				Quantity:  1,
				Price:     10.99,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				UserID:    userID,
				ProductID: 2,
				Quantity:  2,
				Price:     20.99,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		for _, item := range items {
			_, err := repo.Create(ctx, item)
			require.NoError(t, err)
		}

		// Verify items were created
		userItems, err := repo.GetByUser(ctx, userID)
		require.NoError(t, err)
		require.Len(t, userItems, 2)

		// Delete all items
		err = repo.DeleteAllByUser(ctx, userID)
		require.NoError(t, err)

		// Verify all items were deleted
		userItems, err = repo.GetByUser(ctx, userID)
		require.NoError(t, err)
		require.Empty(t, userItems)
	})
}
