BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS transactions
(
    id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_user UUID REFERENCES users (id) ON DELETE CASCADE,
    to_user   UUID REFERENCES users (id) ON DELETE CASCADE,
    amount    INT              DEFAULT 0
);

COMMIT;