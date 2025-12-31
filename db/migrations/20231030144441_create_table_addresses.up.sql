create table if not exists addresses
(
    id          bigserial    not null,
    contact_id  bigint       not null,
    street      varchar(255),
    city        varchar(255),
    province    varchar(255),
    postal_code varchar(10),
    country     varchar(100),
    created_at  timestamptz  not null default now(),
    updated_at  timestamptz  not null default now(),
    primary key (id),
    constraint fk_addresses_contact_id foreign key (contact_id) references contacts (id)
);
