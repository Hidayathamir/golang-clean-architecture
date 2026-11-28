-- +migrate Up
create index idx_chats_recipient_id_read_at
    on chats (recipient_id, read_at);

-- +migrate Down
drop index idx_chats_recipient_id_read_at;