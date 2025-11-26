-- +goose Up
-- +goose StatementBegin
CREATE TABLE delivery (
                          id                  BIGSERIAL PRIMARY KEY,
                          courier_id          BIGINT NOT NULL REFERENCES couriers(id),
                          order_id            VARCHAR(255) NOT NULL UNIQUE,
                          assigned_at         TIMESTAMP NOT NULL DEFAULT NOW(),
                          deadline            TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE delivery;
-- +goose StatementEnd
