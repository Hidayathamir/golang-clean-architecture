-- +migrate Up
create index idx_user_stats_deleted_at
on user_stats (deleted_at);
-- +migrate Down
drop index idx_user_stats_deleted_at;