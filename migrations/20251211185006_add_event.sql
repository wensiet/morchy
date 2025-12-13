-- +goose Up
-- +goose StatementBegin
CREATE TABLE event (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    source_id VARCHAR(36) NOT NULL,
    node_id VARCHAR(36) NOT NULL,
    payload JSONB,
    produced_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event;
-- +goose StatementEnd
