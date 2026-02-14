-- +migrate Up
create unique index idx_users_username_active 
on users (username) 
where (deleted_at is null);
-- +migrate Down
drop index idx_users_username_active;