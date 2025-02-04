-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    status,
    payment_status,
    total_amount,
    shipping_address,
    shipping_city,
    shipping_country,
    shipping_zip,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: GetOrderItems :many
SELECT * FROM order_items WHERE order_id = $1;

-- name: GetOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountOrdersByUserID :one
SELECT COUNT(*) FROM orders WHERE user_id = $1;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2,
    updated_at = $3
WHERE id = $1;

-- name: UpdatePaymentStatus :exec
UPDATE orders
SET payment_status = $2,
    updated_at = $3
WHERE id = $1;
