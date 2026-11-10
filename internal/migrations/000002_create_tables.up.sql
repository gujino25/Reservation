CREATE TABLE IF NOT EXISTS tables ( 
    id BIGSERIAL PRIMARY KEY,
    restaurant_id BIGINT NOT NULL REFERENCES restaurants(id),
    number INT NOT NULL,
    capacity INT NOT NULL,
    UNIQUE (restaurant_id, number)
);