-- name: CreateProduct :one
insert into products(id, created_at, updated_at, name, url, user_id)
values ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: FindProductsByUser :many
select * from products
where user_id = $1
order by updated_at desc;

-- name: GetProductsToFetch :many
select * from products
order by last_fetched_at asc nulls first
limit $1;

-- name: UpdateBasePrice :one
update products
set last_fetched_at = NOW(),
updated_at = NOW(),
base_price = $1
where id = $2
RETURNING *;

-- name: UpdateCurrentPrice :one
update products
set last_fetched_at = NOW(),
updated_at = NOW(),
current_Price = $1
where id = $2
RETURNING *;

-- name: UpdateProductPriority :one
update products
set priority = $1
where id = $2
RETURNING *;

-- name: GetProductById :one
select * from products
where id = $1;

-- name: DeleteProduct :exec
delete from products
where id = $1;
