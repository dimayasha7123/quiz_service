-- +goose Up
-- +goose StatementBegin
alter table question drop column quiz_id;
create table quiz_has_question
(
    quiz_id bigint references quiz,
    question_id bigint references question,
    primary key (quiz_id, question_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table question add column quiz_id bigint references quiz;
drop table if exists quiz_has_question;
-- +goose StatementEnd
