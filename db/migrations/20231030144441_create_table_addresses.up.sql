create table addresses
(
    id          bigserial   primary key,
    contact_id  bigint      not null,
    street      varchar     null,
    city        varchar     null,
    province    varchar     null,
    postal_code varchar     null,
    country     varchar     null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now(),
    deleted_at  timestamptz null,
    constraint fk_addresses_contact_id foreign key (contact_id) references contacts (id)  on delete cascade
);
