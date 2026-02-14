-- +migrate Up
alter table follows add constraint 
fk_follows_following_id foreign key (following_id) references users (id) on delete cascade;

-- +migrate Down
alter table follows drop constraint fk_follows_following_id;
