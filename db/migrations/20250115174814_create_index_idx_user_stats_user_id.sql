-- +migrate Up
create index idx_user_stats_user_id 
on user_stats (user_id);

-- +migrate Down
drop index if exists idx_user_stats_user_id;
