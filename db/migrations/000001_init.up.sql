CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Function to generate account ID with id and 9 random digits
CREATE OR REPLACE FUNCTION generate_account_id(id TEXT) RETURNS TEXT AS
$$
BEGIN
    RETURN upper(id) || lpad(cast(floor(random() * 1000000000) AS text), 9, '0');
END;
$$ LANGUAGE plpgsql;

-- Function to generate transaction ID with id and cuid
CREATE OR REPLACE FUNCTION generate_transaction_id() RETURNS TEXT AS
$$
DECLARE
    cuid TEXT;
BEGIN
    SELECT concat_ws(
                   '-',
                   'TX',
                   upper(to_hex(trunc(extract(epoch from clock_timestamp()) * 1000)::bigint)),
                   upper(lpad(to_hex(floor(random() * power(2, 16))::bigint), 4, '0'))
           )
    INTO cuid;
    RETURN cuid;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS audit
(
    id         UUID      DEFAULT uuid_generate_v4(),
    operation  TEXT                    NOT NULL,
    table_name TEXT                    NOT NULL,
    record_id  TEXT                    NOT NULL,
    actor      TEXT                    NOT NULL,
    actor_id   TEXT,
    remarks    TEXT,
    old_record JSONB,
    new_record JSONB,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    PRIMARY KEY (id, created_at),
    CONSTRAINT check_audit_valid_operation CHECK (operation IN ('CREATE', 'UPDATE', 'DELETE'))
) PARTITION BY RANGE (created_at);

CREATE TABLE IF NOT EXISTS audit_part_2024_2 PARTITION OF audit
    FOR VALUES FROM ('2024-07-01') TO ('2024-12-31');
CREATE TABLE IF NOT EXISTS audit_part_2025_1 PARTITION OF audit
    FOR VALUES FROM ('2025-01-01') TO ('2025-06-30');
CREATE TABLE IF NOT EXISTS audit_part_2025_2 PARTITION OF audit
    FOR VALUES FROM ('2025-07-01') TO ('2025-12-31');
CREATE TABLE IF NOT EXISTS audit_part_2026_1 PARTITION OF audit
    FOR VALUES FROM ('2026-01-01') TO ('2026-06-30');
CREATE TABLE IF NOT EXISTS audit_part_2026_2 PARTITION OF audit
    FOR VALUES FROM ('2026-07-01') TO ('2026-12-31');

CREATE INDEX IF NOT EXISTS audit_operation_at_idx ON audit (created_at);
CREATE INDEX IF NOT EXISTS audit_table_name_idx ON audit (table_name);
CREATE INDEX IF NOT EXISTS audit_record_id_idx ON audit (record_id);
CREATE INDEX IF NOT EXISTS audit_actor_idx ON audit (actor);
CREATE INDEX IF NOT EXISTS audit_actor_id_idx ON audit (actor_id);


CREATE TABLE IF NOT EXISTS wallets
(
    id                  TEXT PRIMARY KEY CHECK (id ~ '^[a-z]+(?:[-_][a-z]+)*$'),
    name                TEXT UNIQUE NOT NULL,
    description         TEXT,
    currency            TEXT        NOT NULL CHECK ( currency ~ '^[a-zA-Z]+$'),
    points_expire_after INTERVAL SECOND(0),                     -- NULL means points never expire
    limit_per_user      BIGINT CHECK (limit_per_user >= 0),     -- NULL means no limit
    limit_global        BIGINT CHECK (limit_global >= 0),       -- NULL means no limit
    minimum_withdrawal  BIGINT CHECK (minimum_withdrawal >= 0), -- NULL means no minimum
    is_monetary         BOOLEAN     NOT NULL DEFAULT FALSE,
    is_active           BOOLEAN     NOT NULL DEFAULT TRUE,
    is_archived         BOOLEAN     NOT NULL DEFAULT FALSE,
    updated_at          TIMESTAMP            DEFAULT NOW() NOT NULL,
    created_at          TIMESTAMP            DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS tiers
(
    id          TEXT PRIMARY KEY CHECK (id ~ '^[a-z]+(?:[-_][a-z]+)*$'),
    name        TEXT                    NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id         TEXT PRIMARY KEY,
    tier_id    TEXT                    REFERENCES tiers (id) ON DELETE SET NULL,
    is_active  BOOLEAN                 NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS exchange_rates
(
    id             SERIAL PRIMARY KEY,
    from_wallet_id TEXT REFERENCES wallets (id) ON DELETE CASCADE NOT NULL,
    to_wallet_id   TEXT REFERENCES wallets (id) ON DELETE CASCADE NOT NULL,
    tier_id        TEXT REFERENCES tiers (id) ON DELETE CASCADE, -- NULL means default rate
    exchange_rate  numeric                                        NOT NULL CHECK (exchange_rate > 0),
    minimum_amount BIGINT CHECK (minimum_amount >= 0),           -- NULL means no minimum
    created_at     TIMESTAMP DEFAULT NOW()                        NOT NULL,
    updated_at     TIMESTAMP DEFAULT NOW()                        NOT NULL,
    CONSTRAINT unique_from_to_tier UNIQUE (from_wallet_id, to_wallet_id, tier_id),
    CONSTRAINT check_different_wallets CHECK (from_wallet_id <> to_wallet_id)
);

CREATE TABLE IF NOT EXISTS triggers
(
    id         SERIAL PRIMARY KEY,
    name       TEXT                    NOT NULL,
    slug       TEXT UNIQUE             NOT NULL CHECK (slug ~ '^[a-z]+(?:[-_][a-z]+)*$'),
    properties JSONB,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS programs
(
    id             SERIAL PRIMARY KEY,
    name           TEXT                                              NOT NULL,
    wallet_id      TEXT REFERENCES wallets (id) ON DELETE CASCADE    NOT NULL,
    trigger_slug   TEXT REFERENCES triggers (slug) ON DELETE CASCADE NOT NULL,
    condition      JSONB                                             NOT NULL,
    effect         JSONB                                             NOT NULL,
    valid_from     TIMESTAMP DEFAULT NOW()                           NOT NULL,
    valid_until    TIMESTAMP, -- NULL means no expiry
    is_active      BOOLEAN   DEFAULT TRUE                            NOT NULL,
    limit_per_user INT CHECK (limit_per_user >= 0),
    limit_global   INT CHECK (limit_global >= 0),
    created_at     TIMESTAMP DEFAULT NOW()                           NOT NULL,
    updated_at     TIMESTAMP DEFAULT NOW()                           NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts
(
    id         TEXT PRIMARY KEY CHECK (length(id) > 9),
    wallet_id  TEXT REFERENCES wallets (id) ON DELETE CASCADE NOT NULL,
    user_id    TEXT REFERENCES users (id) ON DELETE CASCADE   NOT NULL,
    balance    BIGINT    DEFAULT 0 CHECK (balance >= 0),
    version    BIGINT    DEFAULT 0                            NOT NULL CHECK (version >= 0),
    is_active  BOOLEAN   DEFAULT TRUE                         NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()                        NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()                        NOT NULL,
    CONSTRAINT unique_wallet_user UNIQUE (wallet_id, user_id)
);

CREATE TABLE IF NOT EXISTS transactions
(
    id               TEXT DEFAULT generate_transaction_id(),
    type             TEXT                                            NOT NULL,
    wallet_id        TEXT REFERENCES wallets (id) ON DELETE CASCADE  NOT NULL,
    account_id       TEXT REFERENCES accounts (id) ON DELETE CASCADE NOT NULL,
    reason           TEXT                                            NOT NULL,
    metadata         JSONB,
    program_id       INT                                             REFERENCES programs (id) ON DELETE SET NULL,
    amount           BIGINT                                          NOT NULL CHECK ( amount >= 0 ),
    available_amount BIGINT           DEFAULT 0                      NOT NULL CHECK (available_amount BETWEEN amount AND 0), -- amount after expiry or debit
    expire_at        TIMESTAMP,
    version          BIGINT           DEFAULT 0                      NOT NULL CHECK (version >= 0),
    created_at       TIMESTAMP        DEFAULT NOW()                  NOT NULL,
    PRIMARY KEY (id, wallet_id),
    CONSTRAINT check_transaction_type CHECK (type IN ('CREDIT', 'DEBIT')),
    CONSTRAINT check_transaction_reason CHECK (reason IN
                                               ('REWARD', 'PURCHASE', 'REDEEM', 'PENALTY', 'EXPIRED', 'EXCHANGE',
                                                'WITHDRAWAL', 'DEPOSIT')),
    CONSTRAINT check_transaction_integrity CHECK (reason IN ('REWARD', 'DEPOSIT') AND type = 'CREDIT' OR
                                                  reason IN ('PURCHASE', 'REDEEM', 'PENALTY', 'EXPIRED',
                                                             'WITHDRAWAL') AND type = 'DEBIT')
) PARTITION BY HASH (wallet_id);

CREATE TABLE IF NOT EXISTS transactions_part_1 PARTITION OF transactions
    FOR VALUES WITH (MODULUS 5, REMAINDER 0);
CREATE TABLE IF NOT EXISTS transactions_part_2 PARTITION OF transactions
    FOR VALUES WITH (MODULUS 5, REMAINDER 1);
CREATE TABLE IF NOT EXISTS transactions_part_3 PARTITION OF transactions
    FOR VALUES WITH (MODULUS 5, REMAINDER 2);
CREATE TABLE IF NOT EXISTS transactions_part_4 PARTITION OF transactions
    FOR VALUES WITH (MODULUS 5, REMAINDER 3);
CREATE TABLE IF NOT EXISTS transactions_part_5 PARTITION OF transactions
    FOR VALUES WITH (MODULUS 5, REMAINDER 4);

CREATE INDEX IF NOT EXISTS transactions_wallet_id_idx ON transactions (wallet_id);
CREATE INDEX IF NOT EXISTS transactions_account_id_idx ON transactions (account_id);
CREATE INDEX IF NOT EXISTS transactions_program_id_idx ON transactions (program_id);
CREATE INDEX IF NOT EXISTS transactions_expire_at_idx ON transactions (expire_at);
CREATE INDEX IF NOT EXISTS transactions_created_at_idx ON transactions (created_at);
CREATE INDEX IF NOT EXISTS transactions_type_idx ON transactions (type);