-- +goose Up
-- +goose StatementBegin
ALTER TABLE workload ADD COLUMN container JSONB NOT NULL DEFAULT '{}'::jsonb;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workload DROP COLUMN container;
-- +goose StatementEnd
