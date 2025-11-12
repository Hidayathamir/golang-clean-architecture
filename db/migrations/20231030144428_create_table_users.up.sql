create table if not exists users
(
    id         bigserial    not null,
    username   varchar(100) not null,
    name       varchar(100) not null,
    password   varchar(100) not null,
    token      varchar(100) null,
    created_at bigint       not null,
    updated_at bigint       not null,
    primary key (id),
    constraint users_username_key unique (username)
);
