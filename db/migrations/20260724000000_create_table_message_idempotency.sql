-- +migrate Up
create table message_idempotency
(
    idempotency_key varchar(36)  primary key,
    topic           varchar(255) not null,
    partition       integer      not null,
    record_offset   bigint       not null,
    processed_at    timestamptz  not null default now()
);

create index idx_message_idempotency_processed_at on message_idempotency (processed_at);

-- +migrate Down
drop table message_idempotency;
