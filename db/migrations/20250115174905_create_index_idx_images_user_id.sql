-- +migrate Up
create index idx_images_user_id 
on images (user_id);

-- +migrate Down
drop index if exists idx_images_user_id;
