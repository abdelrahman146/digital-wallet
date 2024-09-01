-- function to set updated_at when the row is updated
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TYPE transaction_type AS ENUM ('DEPOSIT', 'WITHDRAW', 'PURCHASE', 'REFUND', 'TRANSFER_IN', 'TRANSFER_OUT');
CREATE TYPE initiator_type AS ENUM ('SYSTEM', 'USER', 'BACKOFFICE');
CREATE TYPE reference_type AS ENUM ('ORDER', 'PAYMENT_TRANSACTION', 'POINTS', 'TRANSFER');

CREATE TABLE wallets
(
    id         uuid PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    user_id    text UNIQUE              NOT NULL,                                  -- user_id is the owner of the wallet
    balance    numeric(18, 2)           NOT NULL DEFAULT 0 CHECK ( balance >= 0 ), -- amount is the balance of the wallet
    version    bigint                   NOT NULL DEFAULT 0,                        -- version is the version of the wallet
    created_at timestamp with time zone NOT NULL DEFAULT now(),                    -- created_at is the time when the wallet is created
    updated_at timestamp with time zone NOT NULL DEFAULT now()                     -- updated_at is the time when the wallet is updated
);

-- Update updated_at when wallet is updated
CREATE TRIGGER set_wallet_updated_at
    BEFORE UPDATE
    ON wallets
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


CREATE TABLE transactions
(
    id               uuid                      DEFAULT uuid_generate_v4(),
    wallet_id        uuid REFERENCES wallets (id),                                     -- wallet_id is the owner of the transaction
    type             transaction_type NOT NULL,                                       -- type is the type of the transaction
    amount           numeric(18, 2)   NOT NULL,                                       -- amount is the amount of the transaction
    reference_id     text,                                                            -- reference_id is the reference of the transaction (e.g. order_id, bank_transaction_id, transaction_id)
    reference_type   reference_type,                                                  -- reference_type is the type of the reference (e.g. order, bank_transaction, transaction)
    initiated_by     initiator_type   NOT NULL,                                       -- initiated_by is the initiator of the transaction
    previous_balance numeric(18, 2)   NOT NULL CHECK (amount >= 0),                   -- previous_balance is the balance before the transaction
    new_balance      numeric(18, 2)   NOT NULL CHECK (amount >= 0),                   -- new_balance is the balance after the transaction
    created_at       timestamp with time zone  DEFAULT NOW(),                         -- created_at is the time when the transaction is created
    version          bigint           NOT NULL DEFAULT 0,                             -- version is the version of the wallet after this transaction
    PRIMARY KEY (id, created_at),                                                     -- primary key is the combination of id and created_at
    CONSTRAINT chk_balance_correct CHECK ( previous_balance + amount = new_balance ), -- check if the balance is correct
    CONSTRAINT chk_transaction_correct CHECK (
        (type = 'DEPOSIT' AND amount >= 0 AND reference_type = 'PAYMENT_TRANSACTION' AND reference_id IS NOT NULL) OR
        (type = 'WITHDRAW' AND amount <= 0 AND reference_type = 'PAYMENT_TRANSACTION' AND reference_id IS NOT NULL) OR
        (type = 'PURCHASE' AND amount <= 0 AND reference_type = 'ORDER' AND reference_id IS NOT NULL) OR
        (type = 'REFUND' AND amount >= 0 AND reference_type = 'ORDER' AND reference_id IS NOT NULL) OR
        (type = 'TRANSFER_IN' AND amount >= 0 AND reference_type = 'TRANSFER' AND reference_id IS NOT NULL) OR
        (type = 'TRANSFER_OUT' AND amount <= 0 AND reference_type IS NULL AND reference_id IS NULL)
        )
) PARTITION BY RANGE (created_at);

CREATE INDEX transactions_wallet_id_idx ON transactions (wallet_id);
CREATE INDEX transactions_created_at_idx ON transactions (created_at);

-- Q1 2024: January 1st to March 31st
CREATE TABLE transactions_2024_q1 PARTITION OF transactions
    FOR VALUES FROM ('2024-01-01 00:00:00+00') TO ('2024-04-01 00:00:00+00');
-- Q2 2024: April 1st to June 30th
CREATE TABLE transactions_2024_q2 PARTITION OF transactions
    FOR VALUES FROM ('2024-04-01 00:00:00+00') TO ('2024-07-01 00:00:00+00');
-- Q3 2024: July 1st to September 30th
CREATE TABLE transactions_2024_q3 PARTITION OF transactions
    FOR VALUES FROM ('2024-07-01 00:00:00+00') TO ('2024-10-01 00:00:00+00');
-- Q4 2024: October 1st to December 31st
CREATE TABLE transactions_2024_q4 PARTITION OF transactions
    FOR VALUES FROM ('2024-10-01 00:00:00+00') TO ('2025-01-01 00:00:00+00');
-- Q1 2025: January 1st to March 31st
CREATE TABLE transactions_2025_q1 PARTITION OF transactions
    FOR VALUES FROM ('2025-01-01 00:00:00+00') TO ('2025-04-01 00:00:00+00');
-- Q2 2025: April 1st to June 30th
CREATE TABLE transactions_2025_q2 PARTITION OF transactions
    FOR VALUES FROM ('2025-04-01 00:00:00+00') TO ('2025-07-01 00:00:00+00');
-- Q3 2025: July 1st to September 30th
CREATE TABLE transactions_2025_q3 PARTITION OF transactions
    FOR VALUES FROM ('2025-07-01 00:00:00+00') TO ('2025-10-01 00:00:00+00');
-- Q4 2025: October 1st to December 31st
CREATE TABLE transactions_2025_q4 PARTITION OF transactions
    FOR VALUES FROM ('2025-10-01 00:00:00+00') TO ('2026-01-01 00:00:00+00');