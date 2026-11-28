-- +migrate Up
create table notifications (
    id bigserial primary key,
    user_id bigint not null,
    actor_id bigint not null,
    type varchar(32) not null,
    image_id bigint not null default 0,
    is_read boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    deleted_at timestamptz
);

create index idx_notifications_user_id_created_at
    on notifications (user_id, created_at desc, id desc);

create index idx_notifications_user_id_is_read
    on notifications (user_id, is_read);

create index idx_notifications_deleted_at
    on notifications (deleted_at);

-- +migrate Down
drop table notifications;