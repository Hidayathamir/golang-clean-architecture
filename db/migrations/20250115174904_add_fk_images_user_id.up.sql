alter table images add constraint 
fk_images_user_id foreign key (user_id) references users (id) on delete cascade;
