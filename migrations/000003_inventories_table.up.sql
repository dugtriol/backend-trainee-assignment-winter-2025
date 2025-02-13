BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS inventories
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID REFERENCES users (id) ON DELETE CASCADE,
    merch_id    INT REFERENCES merch (id) ON DELETE CASCADE,
    quantity    INT              DEFAULT 1
);

COMMIT;