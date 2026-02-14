-- +migrate Up
alter table likes add constraint 
fk_likes_image_id foreign key (image_id) references images (id) on delete cascade;

-- +migrate Down
alter table likes drop constraint fk_likes_image_id;
