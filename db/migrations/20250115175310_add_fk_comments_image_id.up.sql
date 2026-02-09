alter table comments add constraint 
fk_comments_image_id foreign key (image_id) references images (id) on delete cascade;
