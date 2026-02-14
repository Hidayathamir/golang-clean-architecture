-- +migrate Up
create index idx_likes_image_id 
on likes (image_id);

-- +migrate Down
drop index if exists idx_likes_image_id;
