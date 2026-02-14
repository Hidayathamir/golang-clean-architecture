-- +migrate Up
create index idx_comments_user_id 
on comments (user_id);

-- +migrate Down
drop index if exists idx_comments_user_id;
