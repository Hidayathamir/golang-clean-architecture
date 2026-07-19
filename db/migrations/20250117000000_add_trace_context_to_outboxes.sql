-- +migrate Up
alter table outboxes
    add column trace_context text null;

-- +migrate Down
alter table outboxes
    drop column trace_context;
