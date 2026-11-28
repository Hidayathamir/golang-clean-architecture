-- +migrate Up
create index idx_chats_deleted_at
    on chats (deleted_at);

-- +migrate Down
drop index idx_chats_deleted_at;