-- +goose Up
alter table products
add column last_fetched_at timestamp;

-- +goose Down
alter table products
drop column last_fetched_at;