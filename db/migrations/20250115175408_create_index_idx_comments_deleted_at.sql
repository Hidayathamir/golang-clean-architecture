-- +migrate Up
create index idx_comments_deleted_at
on comments (deleted_at);
-- +migrate Down
drop index idx_comments_deleted_at;