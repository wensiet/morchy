-- +goose Up
-- +goose StatementBegin
BEGIN;
ALTER TABLE spec
    ADD COLUMN image  TEXT NOT NULL,
    ADD COLUMN command TEXT[],
    ADD COLUMN env    JSONB;
ALTER TABLE workload DROP COLUMN container;
COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;
ALTER TABLE workload ADD COLUMN container JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE spec
    DROP COLUMN IF EXISTS image,
    DROP COLUMN IF EXISTS command,
    DROP COLUMN IF EXISTS env;
COMMIT;
-- +goose StatementEnd
