-- +migrate Up
alter table follows add constraint 
fk_follows_follower_id foreign key (follower_id) references users (id) on delete cascade;

-- +migrate Down
alter table follows drop constraint fk_follows_follower_id;
