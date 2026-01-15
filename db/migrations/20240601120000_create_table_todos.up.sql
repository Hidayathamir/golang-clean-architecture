create table todos
(
    id           bigserial      primary key,
    user_id      bigint         not null,
    title        varchar        not null,
    description  text           null,
    is_completed boolean        not null default false,
    completed_at timestamptz    null,
    created_at   timestamptz    not null default now(),
    updated_at   timestamptz    not null default now(),
    deleted_at   timestamptz    null,
    constraint fk_todos_user_id foreign key (user_id) references users (id) on delete cascade
);
