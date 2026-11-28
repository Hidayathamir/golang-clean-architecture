-- +migrate Up
-- De-duplication is per consumer group: one event is legitimately processed by
-- several groups on the same topic (e.g. notify-user and batch-stats on
-- user.followed), so scoping the key by group is required.
alter table message_idempotency drop constraint message_idempotency_pkey;
alter table message_idempotency add column consumer_group varchar(255) not null default '';
alter table message_idempotency add primary key (idempotency_key, consumer_group);

-- +migrate Down
alter table message_idempotency drop constraint message_idempotency_pkey;
alter table message_idempotency drop column consumer_group;
alter table message_idempotency add primary key (idempotency_key);