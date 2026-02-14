-- +migrate Up
create table comments
(
    id          bigserial   primary key,
    user_id     bigint      not null,
    image_id    bigint      not null,
    comment     varchar     not null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now(),
    deleted_at  timestamptz null
);

-- +migrate Down
drop table comments;
