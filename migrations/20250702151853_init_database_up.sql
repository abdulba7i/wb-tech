-- +goose Up
CREATE TABLE IF NOT EXISTS delivery (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS payment (
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(255) NOT NULL,
    request_id VARCHAR(255) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(255) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(255) NOT NULL,
    delivery_id INTEGER NOT NULL REFERENCES delivery(id) ON DELETE CASCADE,
    payment_id INTEGER NOT NULL REFERENCES payment(id) ON DELETE CASCADE,
    locale VARCHAR(255) NOT NULL,
    internal_signature VARCHAR(255) NOT NULL,
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255) NOT NULL,
    shardkey VARCHAR(255) NOT NULL,
    sm_id INTEGER NOT NULL,
    oof_shard VARCHAR(255) NOT NULL,
    date_created TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(255) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL
);

-- Индексы (опционально, но полезны)
CREATE INDEX IF NOT EXISTS idx_orders_delivery_id ON orders(delivery_id);
CREATE INDEX IF NOT EXISTS idx_orders_payment_id ON orders(payment_id);
CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);
CREATE INDEX IF NOT EXISTS idx_items_rid ON items(rid);

-- +goose Down
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS delivery;
