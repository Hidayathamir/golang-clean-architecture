-- +migrate Up
create index idx_follows_following_id 
on follows (following_id);

-- +migrate Down
drop index if exists idx_follows_following_id;
