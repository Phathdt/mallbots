package repositories

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"mallbots/modules/product/domain/interfaces"

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

func TestGetProduct(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	product, err := repo.GetProduct(context.Background(), 1)

	require.NoError(t, err)
	require.NotNil(t, product)
	require.Equal(t, "iPhone 15 Pro", product.Name)
	require.Equal(t, 999.99, product.Price)
}

func TestGetProducts(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	testCases := []struct {
		name         string
		filter       *interfaces.ProductFilter
		paging       *core.Paging
		expectedLen  int
		expectedName string
	}{
		{
			name: "Filter by price range",
			filter: &interfaces.ProductFilter{
				MinPrice: &[]float64{900.0}[0],
				MaxPrice: &[]float64{1000.0}[0],
			},
			paging: &core.Paging{
				Page:  1,
				Limit: 10,
			},
			expectedLen:  1,
			expectedName: "iPhone 15 Pro",
		},
		{
			name: "Filter by category",
			filter: &interfaces.ProductFilter{
				Category: &[]int32{1}[0],
			},
			paging: &core.Paging{
				Page:  1,
				Limit: 10,
			},
			expectedLen: 3,
		},
		{
			name: "Search by name",
			filter: &interfaces.ProductFilter{
				Search: "iPhone",
			},
			paging: &core.Paging{
				Page:  1,
				Limit: 10,
			},
			expectedLen:  1,
			expectedName: "iPhone 15 Pro",
		},
		{
			name: "Sort by price ascending",
			filter: &interfaces.ProductFilter{
				SortBy: "price_asc",
			},
			paging: &core.Paging{
				Page:  1,
				Limit: 5,
			},
			expectedLen: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			products, err := repo.GetProducts(context.Background(), tc.filter, tc.paging)

			for _, product := range products {
				fmt.Println(product.Name)
				fmt.Println(product.ID)
			}
			require.NoError(t, err)
			require.Len(t, products, tc.expectedLen)

			if tc.expectedName != "" {
				require.Equal(t, tc.expectedName, products[0].Name)
			}

			if tc.filter.SortBy == "price_asc" {
				for i := 0; i < len(products)-1; i++ {
					require.LessOrEqual(t, products[i].Price, products[i+1].Price)
				}
			}
		})
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	product, err := repo.GetProduct(context.Background(), 999)

	require.Error(t, err)
	require.Nil(t, product)
}

func TestGetProducts_Pagination(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := NewProductRepository(db)

	filter := &interfaces.ProductFilter{}
	paging := &core.Paging{
		Page:  1,
		Limit: 5,
	}

	products, err := repo.GetProducts(context.Background(), filter, paging)

	require.NoError(t, err)
	require.Len(t, products, 5)
	require.Equal(t, int64(23), paging.Total)

	paging.Page = 2
	productsPage2, err := repo.GetProducts(context.Background(), filter, paging)
	require.NoError(t, err)
	require.Len(t, productsPage2, 5)

	require.NotEqual(t, products[0].ID, productsPage2[0].ID)
}
