-- +goose Up
ALTER TABLE accounts RENAME COLUMN type_accpunt_id TO type_account_id;

-- +goose Down
ALTER TABLE accounts RENAME COLUMN type_account_id TO type_accpunt_id;
