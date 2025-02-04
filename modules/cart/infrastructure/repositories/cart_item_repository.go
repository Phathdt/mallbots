package repositories

import (
	"context"
	"mallbots/modules/cart/domain/entities"
	"mallbots/modules/cart/domain/interfaces"
	"mallbots/modules/cart/infrastructure/query/gen"

	"github.com/jackc/pgx/v5/pgxpool"
)

type cartItemRepository struct {
	db *pgxpool.Pool
}

func NewCartItemRepository(db *pgxpool.Pool) interfaces.CartItemRepository {
	return &cartItemRepository{db: db}
}

func (r *cartItemRepository) Create(ctx context.Context, item *entities.CartItem) (*entities.CartItem, error) {
	queries := gen.New(r.db)

	dbItem, err := queries.CreateCartItem(ctx, gen.CreateCartItemParams{
		UserID:    item.UserID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Price:     item.Price,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}

	return &entities.CartItem{
		ID:        dbItem.ID,
		UserID:    dbItem.UserID,
		ProductID: dbItem.ProductID,
		Quantity:  dbItem.Quantity,
		Price:     dbItem.Price,
		CreatedAt: dbItem.CreatedAt,
		UpdatedAt: dbItem.UpdatedAt,
	}, nil
}

func (r *cartItemRepository) Update(ctx context.Context, item *entities.CartItem) error {
	queries := gen.New(r.db)

	return queries.UpdateCartItem(ctx, gen.UpdateCartItemParams{
		UserID:    item.UserID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		UpdatedAt: item.UpdatedAt,
	})
}

func (r *cartItemRepository) Delete(ctx context.Context, userID, productID int32) error {
	queries := gen.New(r.db)

	return queries.DeleteCartItem(ctx, gen.DeleteCartItemParams{
		UserID:    userID,
		ProductID: productID,
	})
}

func (r *cartItemRepository) GetByUserAndProduct(ctx context.Context, userID, productID int32) (*entities.CartItem, error) {
	queries := gen.New(r.db)

	dbItem, err := queries.GetCartItem(ctx, gen.GetCartItemParams{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		return nil, err
	}

	return &entities.CartItem{
		ID:        dbItem.ID,
		UserID:    dbItem.UserID,
		ProductID: dbItem.ProductID,
		Quantity:  dbItem.Quantity,
		Price:     dbItem.Price,
		CreatedAt: dbItem.CreatedAt,
		UpdatedAt: dbItem.UpdatedAt,
	}, nil
}

func (r *cartItemRepository) GetByUser(ctx context.Context, userID int32) ([]*entities.CartItem, error) {
	queries := gen.New(r.db)

	dbItems, err := queries.GetCartItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	items := make([]*entities.CartItem, len(dbItems))
	for i, dbItem := range dbItems {
		items[i] = &entities.CartItem{
			ID:        dbItem.ID,
			UserID:    dbItem.UserID,
			ProductID: dbItem.ProductID,
			Quantity:  dbItem.Quantity,
			Price:     dbItem.Price,
			CreatedAt: dbItem.CreatedAt,
			UpdatedAt: dbItem.UpdatedAt,
		}
	}

	return items, nil
}
