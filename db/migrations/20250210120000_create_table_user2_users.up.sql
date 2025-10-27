create table if not exists user2_users
(
    id            varchar(100) not null,
    email         varchar(255) not null,
    password_hash varchar(255) not null,
    display_name  varchar(100) not null,
    created_at    bigint       not null,
    updated_at    bigint       not null,
    primary key (id),
    unique (email)
);
