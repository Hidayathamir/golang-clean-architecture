create table if not exists users
(
    id         bigserial    not null,
    username   varchar(100) not null,
    name       varchar(100) not null,
    password   varchar(100) not null,
    token      varchar(100) null,
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now(),
    primary key (id),
    constraint users_username_key unique (username)
);
