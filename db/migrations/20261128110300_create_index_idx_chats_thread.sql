-- +migrate Up
create index idx_chats_thread
    on chats (sender_id, recipient_id, created_at desc);

-- +migrate Down
drop index idx_chats_thread;