-- +goose Up
alter table products
add column base_price int,
add column current_price int;

-- +goose Down
alter table products
drop column base_price,
drop column current_price,
drop column last_fetched_at;