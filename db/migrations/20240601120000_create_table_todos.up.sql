create table if not exists todos
(
    id           bigserial    not null,
    user_id      bigint       not null,
    title        varchar(200) not null,
    description  text         null,
    is_completed boolean      not null default false,
    completed_at timestamptz  null,
    created_at   timestamptz  not null default now(),
    updated_at   timestamptz  not null default now(),
    deleted_at   timestamptz  null,
    constraint todos_pkey primary key (id),
    constraint todos_user_fk foreign key (user_id) references users (id) on delete cascade
);

create index if not exists idx_todos_user_id on todos (user_id);
create index if not exists idx_todos_completed on todos (is_completed);
