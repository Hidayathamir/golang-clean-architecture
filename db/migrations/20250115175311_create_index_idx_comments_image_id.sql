-- +migrate Up
create index idx_comments_image_id 
on comments (image_id);

-- +migrate Down
drop index if exists idx_comments_image_id;
