-- Wallets Table in Public Schema
CREATE TABLE IF NOT EXISTS wallets
(
    id                  TEXT PRIMARY KEY CHECK (char_length(id) BETWEEN 1 AND 4) CHECK (id ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$'),
    name                TEXT UNIQUE             NOT NULL,
    description         TEXT,
    currency            TEXT                    NOT NULL CHECK (char_length(currency) BETWEEN 1 AND 4),
    points_expire_after INTERVAL, -- NULL means points never expire
    limit_per_user      BIGINT CHECK (limit_per_user >= 0),
    limit_global        BIGINT CHECK (limit_global >= 0),
    is_monetary         BOOLEAN                 NOT NULL DEFAULT FALSE,
    updated_at          TIMESTAMP DEFAULT NOW() NOT NULL,
    created_at          TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS tiers
(
    id         TEXT PRIMARY KEY CHECK (id ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$'),
    name       TEXT                    NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id         TEXT PRIMARY KEY,
    tier_id    TEXT                    REFERENCES tiers (id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS exchange_rates
(
    id             SERIAL PRIMARY KEY,
    from_wallet_id TEXT REFERENCES wallets (id) ON DELETE CASCADE NOT NULL,
    to_wallet_id   TEXT REFERENCES wallets (id) ON DELETE CASCADE NOT NULL,
    tier_id        TEXT REFERENCES tiers (id) ON DELETE CASCADE,
    exchange_rate  numeric                                        NOT NULL CHECK (exchange_rate > 0),
    created_at     TIMESTAMP DEFAULT NOW()                        NOT NULL,
    updated_at     TIMESTAMP DEFAULT NOW()                        NOT NULL,
    CONSTRAINT unique_from_to_tier UNIQUE (from_wallet_id, to_wallet_id, tier_id),
    CONSTRAINT check_different_wallets CHECK (from_wallet_id <> to_wallet_id)
);

CREATE TABLE IF NOT EXISTS triggers
(
    id         SERIAL PRIMARY KEY,
    name       TEXT                    NOT NULL,
    slug       TEXT UNIQUE             NOT NULL CHECK (slug ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$'),
    schema     JSONB                   NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TYPE transaction_type AS ENUM ('CREDIT', 'DEBIT');
CREATE TYPE transaction_actor_type AS ENUM ('SYSTEM', 'USER', 'BACKOFFICE');

-- Function to create a schema and its tables
CREATE OR REPLACE FUNCTION create_wallet_schema(id TEXT) RETURNS VOID AS
$$
DECLARE
    schema_name TEXT := lower(id) || '_wallet';
BEGIN
    -- Create schema
    EXECUTE 'CREATE SCHEMA IF NOT EXISTS ' || schema_name;

    -- Create accounts table
    EXECUTE 'CREATE TABLE IF NOT EXISTS ' || schema_name || '.accounts (
        id TEXT PRIMARY KEY DEFAULT generate_account_id(' || quote_literal(id) || '),
        user_id TEXT REFERENCES public.users(id) ON DELETE CASCADE UNIQUE NOT NULL,
        balance BIGINT DEFAULT 0 CHECK (balance >= 0),
        version BIGINT DEFAULT 0 NOT NULL CHECK (version >= 0),
        updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
        created_at TIMESTAMP DEFAULT NOW() NOT NULL
    )';

    -- Create partitioned transactions table
    EXECUTE 'CREATE TABLE IF NOT EXISTS ' || schema_name || '.transactions (
        id TEXT PRIMARY KEY DEFAULT generate_transaction_id(' || quote_literal(id) || '),
        type transaction_type NOT NULL,
        account_id TEXT REFERENCES ' || schema_name || '.accounts(id) ON DELETE CASCADE NOT NULL,
        actor_type transaction_actor_type NOT NULL,
        actor_id TEXT NOT NULL,
        metadata JSONB,
        program_id INT REFERENCES ' || schema_name || '.programs(id) ON DELETE SET NULL,
        amount BIGINT NOT NULL,
        previous_balance BIGINT NOT NULL,
        new_balance BIGINT NOT NULL,
        version BIGINT DEFAULT 0 NOT NULL CHECK (version >= 0),
        created_at TIMESTAMP DEFAULT NOW() NOT NULL,
        CONSTRAINT check_credit_debit_balance CHECK (
            (type = ''CREDIT'' AND amount + previous_balance = new_balance) OR
            (type = ''DEBIT'' AND previous_balance - amount = new_balance)
        )
    ) PARTITION BY HASH (id)';

    EXECUTE 'CREATE INDEX IF NOT EXISTS transactions_account_id_idx ON ' || schema_name ||
            '.transactions (account_id)';
    EXECUTE 'CREATE INDEX IF NOT EXISTS transactions_type_idx ON ' || schema_name || '.transactions (type)';
    EXECUTE 'CREATE INDEX IF NOT EXISTS transactions_actor_type_idx ON ' || schema_name ||
            '.transactions (transaction_actor_type)';

    -- Create transactions partitions (standard strategy - hash modulus)
    EXECUTE 'CREATE TABLE IF NOT EXISTS ' || schema_name || '.transactions_part_1 PARTITION OF ' || schema_name || '.transactions
        FOR VALUES WITH (MODULUS 4, REMAINDER 0)';

    EXECUTE 'CREATE TABLE IF NOT EXISTS ' || schema_name || '.transactions_part_2 PARTITION OF ' || schema_name || '.transactions
        FOR VALUES WITH (MODULUS 4, REMAINDER 1)';

    EXECUTE 'CREATE TABLE IF NOT EXISTS ' || schema_name || '.transactions_part_3 PARTITION OF ' || schema_name || '.transactions
        FOR VALUES WITH (MODULUS 4, REMAINDER 2)';

    EXECUTE 'CREATE TABLE IF NOT EXISTS ' || schema_name || '.transactions_part_4 PARTITION OF ' || schema_name || '.transactions
        FOR VALUES WITH (MODULUS 4, REMAINDER 3)';

    EXECUTE 'CREATE TABLE IF NOT EXISTS ' || schema_name || '.programs (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        trigger_slug TEXT REFERENCES public.triggers(slug) ON DELETE CASCADE NOT NULL,
        condition JSONB NOT NULL,
        effect JSONB NOT NULL,
        limit_per_user INT CHECK (limit_per_user >= 0),
        limit_global INT CHECK (limit_global >= 0),
        created_at TIMESTAMP DEFAULT NOW() NOT NULL
    )';

END;
$$ LANGUAGE plpgsql;

-- Function to generate account ID with id and 9 random digits
CREATE OR REPLACE FUNCTION generate_account_id(id TEXT) RETURNS TEXT AS
$$
BEGIN
    RETURN upper(id) || lpad(cast(floor(random() * 1000000000) AS text), 9, '0');
END;
$$ LANGUAGE plpgsql;

-- Function to generate transaction ID with id and cuid
CREATE OR REPLACE FUNCTION generate_transaction_id(id TEXT) RETURNS TEXT AS
$$
DECLARE
    cuid TEXT;
BEGIN
    -- Generate cuid (collision-resistant unique identifier)
    SELECT concat_ws(
                   '-',
                   to_hex(trunc(extract(epoch from clock_timestamp()) * 1000)::bigint),
                   lpad(to_hex(floor(random() * 4294967296)::bigint), 8, '0')
           )
    INTO cuid;

    RETURN upper(id) || cuid;
END;
$$ LANGUAGE plpgsql;
