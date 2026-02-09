alter table user_stats add constraint 
fk_user_stats_user_id foreign key (user_id) references users (id) on delete cascade;
