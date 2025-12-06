-- +goose Up
-- +goose StatementBegin
CREATE TABLE workload (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE spec (
    id VARCHAR(36) PRIMARY KEY REFERENCES workload NOT NULL,
    cpu INTEGER NOT NULL,
    ram INTEGER NOT NULL
);

CREATE TABLE lease (
    node_id VARCHAR(36) NOT NULL,
    workload_id VARCHAR(36) UNIQUE REFERENCES workload NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_lease_updated_at
    BEFORE UPDATE ON lease
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_lease_updated_at ON lease (updated_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_lease_updated_at;
DROP TRIGGER trigger_lease_updated_at;
DROP FUNCTION set_updated_at;
DROP TABLE lease;
DROP TABLE spec;
DROP TABLE workload;
-- +goose StatementEnd
