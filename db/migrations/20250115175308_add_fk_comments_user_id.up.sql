alter table comments add constraint 
fk_comments_user_id foreign key (user_id) references users (id) on delete cascade;
