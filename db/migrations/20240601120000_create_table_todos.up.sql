create table if not exists todos
(
    id           varchar(100) not null,
    user_id      bigint       not null,
    title        varchar(200) not null,
    description  text         null,
    is_completed boolean      not null default false,
    completed_at bigint       null,
    created_at   bigint       not null,
    updated_at   bigint       not null,
    constraint todos_pkey primary key (id),
    constraint todos_user_fk foreign key (user_id) references users (id) on delete cascade
);

create index if not exists idx_todos_user_id on todos (user_id);
create index if not exists idx_todos_completed on todos (is_completed);
