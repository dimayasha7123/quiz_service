-- +goose Up
-- +goose StatementBegin
alter table participation
    rename to party;
alter table quiz
    rename column name to title;
alter table response_report
    rename column participation_id to party_id;

create table question
(
    id      bigserial primary key,
    quiz_id bigint references quiz,
    title   varchar not null
);

create table answer
(
    id          bigserial primary key,
    question_id bigint references question,
    title       varchar not null,
    correct     bool default false
);

create table partyQuestion
(
    question_id bigint references question,
    party_id    bigint references question,
    primary key (question_id, party_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table party
    rename to participation;
alter table quiz
    rename column title to name;
alter table response_report
    rename column party_id to participation_id;

drop table if exists answer;
drop table if exists partyQuestion;
drop table if exists question;
-- +goose StatementEnd
