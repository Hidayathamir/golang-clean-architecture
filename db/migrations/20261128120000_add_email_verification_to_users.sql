-- +migrate Up
alter table users
    add column email varchar(255) not null default '',
    add column email_verified_at timestamptz null,
    add column verification_token varchar(255) not null default '';

create unique index idx_users_email on users (email) where email <> '';

-- +migrate Down
drop index idx_users_email;

alter table users
    drop column email,
    drop column email_verified_at,
    drop column verification_token;