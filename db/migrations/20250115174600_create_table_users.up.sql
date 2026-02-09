create table users
(
    id              bigserial   primary key,
    username        varchar     not null,
    password        text        not null,
    name            varchar     not null,
    created_at      timestamptz not null default now(),
    updated_at      timestamptz not null default now(),
    deleted_at      timestamptz null
);
