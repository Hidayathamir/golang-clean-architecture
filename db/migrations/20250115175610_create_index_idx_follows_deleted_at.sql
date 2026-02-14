-- +migrate Up
create index idx_follows_deleted_at
on follows (deleted_at);
-- +migrate Down
drop index idx_follows_deleted_at;