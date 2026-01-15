create table contacts
(
    id         bigserial    primary key,
    first_name varchar      not null,
    last_name  varchar      null,
    email      varchar      null,
    phone      varchar      null,
    user_id    bigint       not null,
    created_at timestamptz  not null default now(),
    updated_at timestamptz  not null default now(),
    deleted_at timestamptz  null,
    constraint fk_contacts_user_id foreign key (user_id) references users (id)  on delete cascade
);
