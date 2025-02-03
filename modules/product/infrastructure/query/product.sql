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

-- name: GetProducts :many
SELECT * FROM products
WHERE
    ($1::text IS NULL OR name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
    AND ($2::int IS NULL OR category_id = $2)
    AND ($3::float IS NULL OR price >= $3)
    AND ($4::float IS NULL OR price <= $4)
ORDER BY
    CASE $5::text
        WHEN 'price_asc' THEN price
        WHEN 'price_desc' THEN price
        WHEN 'latest' THEN extract(epoch from created_at)
        ELSE extract(epoch from created_at)
    END DESC,
    CASE $5::text
        WHEN 'price_asc' THEN id
        ELSE id
    END DESC
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
