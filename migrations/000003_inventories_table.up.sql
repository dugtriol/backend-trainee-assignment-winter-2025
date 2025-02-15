BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS inventories
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID REFERENCES users (id) ON DELETE CASCADE,
    type  VARCHAR(255) REFERENCES merch (name) ON DELETE CASCADE,
    quantity    INT              DEFAULT 1
);

ALTER TABLE inventories ADD CONSTRAINT unique_customer_merch UNIQUE (customer_id, type);

COMMIT;