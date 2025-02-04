-- name: CreateCartItem :one
INSERT INTO cart_items (
    user_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateCartItem :exec
UPDATE cart_items
SET quantity = $3,
    updated_at = $4
WHERE user_id = $1 AND product_id = $2;

-- name: DeleteCartItem :exec
DELETE FROM cart_items
WHERE user_id = $1 AND product_id = $2;

-- name: GetCartItem :one
SELECT * FROM cart_items
WHERE user_id = $1 AND product_id = $2;

-- name: GetCartItems :many
SELECT * FROM cart_items
WHERE user_id = $1
ORDER BY created_at DESC;
