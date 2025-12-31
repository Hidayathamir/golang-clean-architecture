create table if not exists contacts
(
    id         bigserial    not null,
    first_name varchar(100) not null,
    last_name  varchar(100) null,
    email      varchar(100) null,
    phone      varchar(100) null,
    user_id    bigint       not null,
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now(),
    primary key (id),
    constraint fk_contacts_user_id foreign key (user_id) references users (id)
);
