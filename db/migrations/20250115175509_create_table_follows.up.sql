create table follows
(
    id              bigserial   primary key,
    follower_id     bigint      not null,
    following_id    bigint      not null,
    created_at      timestamptz not null default now(),
    updated_at      timestamptz not null default now(),
    deleted_at      timestamptz null
);
