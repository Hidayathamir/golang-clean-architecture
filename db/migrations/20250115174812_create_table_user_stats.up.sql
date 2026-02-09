create table user_stats
(
    id              bigserial   primary key,
    user_id         bigint      not null unique,
    follower_count  int         not null default 0,
    following_count int         not null default 0,
    created_at      timestamptz not null default now(),
    updated_at      timestamptz not null default now(),
    deleted_at      timestamptz null
);
