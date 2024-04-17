-- +goose Up
alter table products
add column priority bool;

update products
set priority = false;

-- +goose Down
alter table products
drop column priority;