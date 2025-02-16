BEGIN;

CREATE TABLE IF NOT EXISTS merch
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(255) NOT NULL UNIQUE,
    price INT      NOT NULL CHECK (price > 0)
);

INSERT INTO merch (name, price)
VALUES ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500) ON CONFLICT (name) DO NOTHING;
COMMIT;