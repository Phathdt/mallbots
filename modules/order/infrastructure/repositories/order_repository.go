package repositories

import (
	"context"
	"mallbots/modules/order/domain/constants"
	"mallbots/modules/order/domain/entities"
	"mallbots/modules/order/domain/interfaces"
	"mallbots/modules/order/infrastructure/query/gen"
	"mallbots/shared/errorx"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phathdt/service-context/core"
)

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) interfaces.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	queries := gen.New(r.db)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)

	dbOrder, err := qtx.CreateOrder(ctx, gen.CreateOrderParams{
		UserID:          order.UserID,
		Status:          order.Status.String(),
		PaymentStatus:   order.PaymentStatus.String(),
		TotalAmount:     order.TotalAmount,
		ShippingAddress: order.ShippingAddress,
		ShippingCity:    order.ShippingCity,
		ShippingCountry: order.ShippingCountry,
		ShippingZip:     order.ShippingZip,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	})
	if err != nil {
		return nil, errorx.ErrCannotCreateOrder
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &entities.Order{
		ID:              dbOrder.ID,
		UserID:          dbOrder.UserID,
		Status:          constants.OrderStatus(dbOrder.Status),
		PaymentStatus:   constants.PaymentStatus(dbOrder.PaymentStatus),
		TotalAmount:     dbOrder.TotalAmount,
		ShippingAddress: dbOrder.ShippingAddress,
		ShippingCity:    dbOrder.ShippingCity,
		ShippingCountry: dbOrder.ShippingCountry,
		ShippingZip:     dbOrder.ShippingZip,
		CreatedAt:       dbOrder.CreatedAt,
		UpdatedAt:       dbOrder.UpdatedAt,
	}, nil
}

func (r *orderRepository) CreateOrderItems(ctx context.Context, orderID int32, items []*entities.OrderItem) error {
	queries := gen.New(r.db)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)

	for _, item := range items {
		_, err := qtx.CreateOrderItem(ctx, gen.CreateOrderItemParams{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
		if err != nil {
			return errorx.ErrCannotCreateOrderItems
		}
	}

	return tx.Commit(ctx)
}

func (r *orderRepository) GetByID(ctx context.Context, id int32) (*entities.Order, error) {
	queries := gen.New(r.db)

	dbOrder, err := queries.GetOrderByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errorx.ErrOrderNotFound
		}
		return nil, err
	}

	// Get order items
	dbItems, err := queries.GetOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}

	var items []*entities.OrderItem
	for _, dbItem := range dbItems {
		items = append(items, &entities.OrderItem{
			ID:        dbItem.ID,
			OrderID:   dbItem.OrderID,
			ProductID: dbItem.ProductID,
			Quantity:  dbItem.Quantity,
			Price:     dbItem.Price,
			CreatedAt: dbItem.CreatedAt,
			UpdatedAt: dbItem.UpdatedAt,
		})
	}

	return &entities.Order{
		ID:              dbOrder.ID,
		UserID:          dbOrder.UserID,
		Status:          constants.OrderStatus(dbOrder.Status),
		PaymentStatus:   constants.PaymentStatus(dbOrder.PaymentStatus),
		TotalAmount:     dbOrder.TotalAmount,
		ShippingAddress: dbOrder.ShippingAddress,
		ShippingCity:    dbOrder.ShippingCity,
		ShippingCountry: dbOrder.ShippingCountry,
		ShippingZip:     dbOrder.ShippingZip,
		CreatedAt:       dbOrder.CreatedAt,
		UpdatedAt:       dbOrder.UpdatedAt,
		Items:           items,
	}, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID int32, paging *core.Paging) ([]*entities.Order, error) {
	queries := gen.New(r.db)

	// Get total count for pagination
	total, err := queries.CountOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	paging.Total = total

	offset := (paging.Page - 1) * paging.Limit

	dbOrders, err := queries.GetOrdersByUserID(ctx, gen.GetOrdersByUserIDParams{
		UserID: userID,
		Limit:  int32(paging.Limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	var orders []*entities.Order
	for _, dbOrder := range dbOrders {
		// Get order items for each order
		dbItems, err := queries.GetOrderItems(ctx, dbOrder.ID)
		if err != nil {
			return nil, err
		}

		var items []*entities.OrderItem
		for _, dbItem := range dbItems {
			items = append(items, &entities.OrderItem{
				ID:        dbItem.ID,
				OrderID:   dbItem.OrderID,
				ProductID: dbItem.ProductID,
				Quantity:  dbItem.Quantity,
				Price:     dbItem.Price,
				CreatedAt: dbItem.CreatedAt,
				UpdatedAt: dbItem.UpdatedAt,
			})
		}

		orders = append(orders, &entities.Order{
			ID:              dbOrder.ID,
			UserID:          dbOrder.UserID,
			Status:          constants.OrderStatus(dbOrder.Status),
			PaymentStatus:   constants.PaymentStatus(dbOrder.PaymentStatus),
			TotalAmount:     dbOrder.TotalAmount,
			ShippingAddress: dbOrder.ShippingAddress,
			ShippingCity:    dbOrder.ShippingCity,
			ShippingCountry: dbOrder.ShippingCountry,
			ShippingZip:     dbOrder.ShippingZip,
			CreatedAt:       dbOrder.CreatedAt,
			UpdatedAt:       dbOrder.UpdatedAt,
			Items:           items,
		})
	}

	return orders, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id int32, status constants.OrderStatus) error {
	queries := gen.New(r.db)

	err := queries.UpdateOrderStatus(ctx, gen.UpdateOrderStatusParams{
		ID:        id,
		Status:    status.String(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return errorx.ErrCannotUpdateOrder
	}

	return nil
}

func (r *orderRepository) UpdatePaymentStatus(ctx context.Context, id int32, status constants.PaymentStatus) error {
	queries := gen.New(r.db)

	err := queries.UpdatePaymentStatus(ctx, gen.UpdatePaymentStatusParams{
		ID:            id,
		PaymentStatus: status.String(),
		UpdatedAt:     time.Now(),
	})
	if err != nil {
		return errorx.ErrCannotUpdateOrder
	}

	return nil
}
