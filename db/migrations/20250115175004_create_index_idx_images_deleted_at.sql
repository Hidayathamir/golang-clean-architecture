-- +migrate Up
create index idx_images_deleted_at
on images (deleted_at);
-- +migrate Down
drop index idx_images_deleted_at;