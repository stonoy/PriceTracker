-- +goose Up
create table products(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text not null,
    url text not null,
    user_id uuid not null
    references users(id)
    on delete cascade,
    unique(url, user_id)
);

-- +goose Down
drop table products;