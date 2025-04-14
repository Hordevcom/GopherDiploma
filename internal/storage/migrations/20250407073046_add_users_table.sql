-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    username TEXT NOT NULL PRIMARY KEY,
    user_password TEXT NOT NULL
);

DO $$ BEGIN
    CREATE TYPE order_status AS ENUM (
        'NEW',
        'PROCESSING',
        'INVALID',
        'PROCESSED'
    );
EXCEPTION
    WHEN duplicate_object THEN NULL;  -- Игнорируем ошибку, если тип уже существует
END $$;

CREATE TABLE IF NOT EXISTS orders (
    number TEXT NOT NULL PRIMARY KEY,
    status order_status NOT NULL DEFAULT 'NEW',
    accrual DOUBLE PRECISION,
    uploaded_at TIMESTAMPTZ NOT NULL,
    username TEXT REFERENCES users(username) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
