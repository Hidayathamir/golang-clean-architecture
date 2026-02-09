alter table follows add constraint 
fk_follows_follower_id foreign key (follower_id) references users (id) on delete cascade;
