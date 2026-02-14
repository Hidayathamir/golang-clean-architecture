-- +migrate Up
alter table user_stats add constraint 
fk_user_stats_user_id foreign key (user_id) references users (id) on delete cascade;

-- +migrate Down
alter table user_stats drop constraint fk_user_stats_user_id;
