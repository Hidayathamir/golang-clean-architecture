create table likes
(
    id          bigserial   primary key,
    user_id     bigint      not null,
    image_id    bigint      not null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now(),
    deleted_at  timestamptz null
);
