-- +goose Up
-- +goose StatementBegin
ALTER TABLE spec
ADD COLUMN host_port INTEGER,
ADD COLUMN container_port INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE spec DROP COLUMN host_port, DROP COLUMN container_port;
-- +goose StatementEnd