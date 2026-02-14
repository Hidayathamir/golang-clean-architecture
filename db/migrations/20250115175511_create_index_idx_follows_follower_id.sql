-- +migrate Up
create index idx_follows_follower_id 
on follows (follower_id);

-- +migrate Down
drop index if exists idx_follows_follower_id;
