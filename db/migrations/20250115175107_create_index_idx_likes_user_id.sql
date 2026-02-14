-- +migrate Up
create index idx_likes_user_id 
on likes (user_id);

-- +migrate Down
drop index if exists idx_likes_user_id;
