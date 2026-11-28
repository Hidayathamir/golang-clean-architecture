-- +migrate Up
alter table chats
    add constraint fk_chats_recipient_id foreign key (recipient_id) references users (id);

-- +migrate Down
alter table chats drop constraint fk_chats_recipient_id;