-- name: ListProducts :many
SELECT * 
FROM products;

-- name: FindProductByID :one
SELECT *
FROM products
WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (customer_id, created_at)
VALUES ($1, NOW())
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_cents)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = $1;

-- name: ListOrderItemsByOrderID :many
SELECT product_id, quantity, price_cents
FROM order_items
WHERE order_id = $1;

-- name: DecrementProductStock :execrows
UPDATE products
SET quantity = quantity - $1
WHERE id = $2
AND quantity >= $1;
