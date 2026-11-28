-- +migrate Up
alter table chats
    add constraint fk_chats_sender_id foreign key (sender_id) references users (id);

-- +migrate Down
alter table chats drop constraint fk_chats_sender_id;