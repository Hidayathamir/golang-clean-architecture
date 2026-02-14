-- +migrate Up
create index idx_likes_deleted_at
on likes (deleted_at);
-- +migrate Down
drop index idx_likes_deleted_at;