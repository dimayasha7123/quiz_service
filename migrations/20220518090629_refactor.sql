-- +goose Up
-- +goose StatementBegin
alter table partyquestion
rename to party_question;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table party_question
    rename to partyquestion;
-- +goose StatementEnd
