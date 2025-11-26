-- +goose Up
-- +goose StatementBegin
ALTER TABLE couriers
    ADD COLUMN total_deliveries INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE couriers
    DROP COLUMN total_deliveries;
-- +goose StatementEnd
