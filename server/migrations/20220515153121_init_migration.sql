-- +goose Up
-- +goose StatementBegin
create table user_account
(
    id   bigserial primary key,
    name varchar unique not null
);

create table quiz
(
    id   bigserial primary key,
    name varchar not null
);

create table participation
(
    id              bigserial primary key,
    user_account_id bigint references user_account,
    quiz_id         bigint references quiz
);

create table response_report
(
    id               bigserial primary key,
    participation_id bigint references participation,
    correct          bool    not null,
    penalty_time     integer not null
);
insert into quiz (name)
values  ('DevOps'),
        ('JavaScript'),
        ('PHP'),
        ('BASH'),
        ('HTML'),
        ('Laravel'),
        ('Docker'),
        ('Linux'),
        ('Kubernetes'),
        ('MySQL');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists response_report;
drop table if exists participation;
drop table if exists user_account;
drop table if exists quiz;
-- +goose StatementEnd
