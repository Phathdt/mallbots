package repositories

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"mallbots/modules/user/domain/entities"

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

func TestUserRepository_Create(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	testCases := []struct {
		name    string
		user    *entities.User
		wantErr bool
	}{
		{
			name: "Create valid user",
			user: &entities.User{
				Email:     "test@example.com",
				Password:  "hashedpassword",
				FullName:  "Test User",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Create user with duplicate email",
			user: &entities.User{
				Email:     "test@example.com",
				Password:  "hashedpassword2",
				FullName:  "Test User 2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.Create(context.Background(), tc.user)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.NotZero(t, user.ID)
				require.Equal(t, tc.user.Email, user.Email)
				require.Equal(t, tc.user.FullName, user.FullName)
			}
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// First create a test user
	testUser := &entities.User{
		Email:     "get@example.com",
		Password:  "hashedpassword",
		FullName:  "Get Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	testCases := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Get existing user",
			email:   "get@example.com",
			wantErr: false,
		},
		{
			name:    "Get non-existing user",
			email:   "nonexistent@example.com",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.GetByEmail(context.Background(), tc.email)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.Equal(t, tc.email, user.Email)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// First create a test user
	testUser := &entities.User{
		Email:     "getid@example.com",
		Password:  "hashedpassword",
		FullName:  "GetID Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := repo.Create(context.Background(), testUser)
	require.NoError(t, err)
	require.NotNil(t, createdUser)

	testCases := []struct {
		name    string
		id      int32
		wantErr bool
	}{
		{
			name:    "Get existing user",
			id:      createdUser.ID,
			wantErr: false,
		},
		{
			name:    "Get non-existing user",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.GetByID(context.Background(), tc.id)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.Equal(t, tc.id, user.ID)
			}
		})
	}
}
