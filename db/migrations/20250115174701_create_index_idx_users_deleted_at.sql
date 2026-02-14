-- +migrate Up
create index idx_users_deleted_at
on users (deleted_at);
-- +migrate Down
drop index idx_users_deleted_at;