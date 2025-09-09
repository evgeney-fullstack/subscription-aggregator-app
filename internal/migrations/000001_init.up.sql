--struct Subscription
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY UNIQUE,
    service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    user_id UUID ,
    start_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finish_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

