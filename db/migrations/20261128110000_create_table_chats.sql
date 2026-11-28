-- +migrate Up
create table chats (
    id bigserial primary key,
    sender_id bigint not null,
    recipient_id bigint not null,
    content text not null,
    read_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

-- +migrate Down
drop table chats;