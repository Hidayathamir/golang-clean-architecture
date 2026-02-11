create table images
(
    id              bigserial   primary key,
    user_id         bigint      not null,
    caption         text        not null,
    url             text        not null,
    like_count      int     not null,
    comment_count   int     not null,
    created_at      timestamptz not null default now(),
    updated_at      timestamptz not null default now(),
    deleted_at      timestamptz null
);
