-- +migrate Up
create table outboxes
(
    id         bigserial    primary key,
    topic      varchar(255) not null,
    key        varchar(255) null,
    payload    bytea        not null,
    status     varchar(50)  not null default 'pending',
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now(),
    deleted_at timestamptz  null
);

create index idx_outboxes_status on outboxes (status);
create index idx_outboxes_created_at on outboxes (created_at);

-- +migrate Down
drop table outboxes;
