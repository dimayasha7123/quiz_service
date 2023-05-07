-- +goose Up
-- +goose StatementBegin
drop table if exists response_report;
alter table party add column completed bool default false;
alter table party add column points integer default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
create table response_report
(
    id               bigserial primary key,
    party_id         bigint references party,
    correct          bool    not null,
    penalty_time     integer not null
);
alter table party drop column completed;
alter table party drop column points;
-- +goose StatementEnd
