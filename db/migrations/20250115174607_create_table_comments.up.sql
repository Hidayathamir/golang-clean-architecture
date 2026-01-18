create table comments
(
    id          bigserial   primary key,
    user_id     bigint      not null,
    image_id    bigint      not null,
    comment     varchar     not null,
    created_at  timestamptz not null default now(),
    updated_at  timestamptz not null default now(),
    deleted_at  timestamptz null,
    constraint fk_comments_user_id foreign key (user_id) references users (id) on delete cascade,
    constraint fk_comments_image_id foreign key (image_id) references images (id) on delete cascade
);
