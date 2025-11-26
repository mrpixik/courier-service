-- +goose Up
-- +goose StatementBegin
ALTER TABLE couriers
ADD transport_type TEXT NOT NULL DEFAULT 'on_foot';  -- on_foot | scooter | car
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE couriers
DROP COLUMN transport_type;
-- +goose StatementEnd
