-- name: CreateProduct :one
INSERT INTO products (
    name,
    description,
    price,
    category_id,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, NOW(), NOW()
) RETURNING *;

-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: CountProducts :one
SELECT COUNT(*) FROM products
WHERE
    (NULLIF(TRIM($1), '') IS NULL OR name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
    AND ($2 = 0 OR category_id = $2)
    AND ($3 = 0 OR price >= $3)
    AND ($4 = 0 OR price <= $4);

-- name: GetProducts :many
SELECT * FROM products
WHERE
    (NULLIF(TRIM($1), '') IS NULL OR name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
    AND ($2 = 0 OR category_id = $2)
    AND ($3 = 0 OR price >= $3)
    AND ($4 = 0 OR price <= $4)
ORDER BY
    CASE $5::text
        WHEN 'price_asc' THEN price
        WHEN 'price_desc' THEN price * -1
        ELSE extract(epoch from created_at) * -1
    END,
    id DESC
LIMIT $6 OFFSET $7;

-- name: GetCategory :one
SELECT * FROM categories WHERE id = $1;

-- name: GetCategories :many
SELECT * FROM categories
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetProductsByCategory :many
SELECT * FROM products
WHERE category_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCategoriesByIds :many
SELECT * FROM categories
WHERE id = ANY($1::int[]);
