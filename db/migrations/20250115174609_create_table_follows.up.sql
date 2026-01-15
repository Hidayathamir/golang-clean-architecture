create table follows
(
    id              bigserial   primary key,
    follower_id     bigint      not null,
    following_id    bigint      not null,
    created_at      timestamptz not null default now(),
    updated_at      timestamptz not null default now(),
    deleted_at      timestamptz null,
    constraint fk_follows_follower_id foreign key (follower_id) references users (id) on delete cascade,
    constraint fk_follows_following_id foreign key (following_id) references users (id) on delete cascade
);
