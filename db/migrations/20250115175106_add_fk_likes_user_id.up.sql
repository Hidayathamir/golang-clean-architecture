alter table likes add constraint 
fk_likes_user_id foreign key (user_id) references users (id) on delete cascade;
