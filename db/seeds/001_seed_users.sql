-- +migrate Up
insert into users (username, password, name)
values ('alice', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Alice'),
       ('bob', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Bob')
on conflict do nothing;

insert into user_stats (user_id, follower_count, following_count)
select id, 0, 0 from users
on conflict (user_id) do nothing;

-- +migrate Down
delete from users where username in ('alice', 'bob');