-- +goose Up
-- +goose StatementBegin
create table user_account
(
    id         integer primary key,
    th_chat_id integer not null,
    username   varchar not null
);

create table quiz
(
    id   integer primary key,
    name varchar not null
);

create table participation
(
    id              integer primary key,
    user_account_id integer references user_account,
    quiz_id         integer references quiz
);

create table response_report
(
    id               integer primary key,
    participation_id integer references participation,
    correct          bool      not null,
    penalty_time     integer not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table response_report;
drop table participation;
drop table user_account;
drop table quiz;
-- +goose StatementEnd
