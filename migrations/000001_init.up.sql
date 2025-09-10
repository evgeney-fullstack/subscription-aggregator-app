CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY UNIQUE,
    service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    user_id UUID,
    start_date DATE NOT NULL,
    finish_date DATE GENERATED ALWAYS AS (start_date + INTERVAL '1 month') STORED
);