-- +goose Up
-- +goose StatementBegin
alter table party_question
    add column quest_order_number integer default 0;
alter table answer
    add order_number integer default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table party_question
    drop column quest_order_number;
alter table answer
    drop column order_number;
-- +goose StatementEnd
