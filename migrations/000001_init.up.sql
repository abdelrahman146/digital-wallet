CREATE TABLE wallet
(
    id         uuid PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    user_id    text UNIQUE              NOT NULL,                                 -- user_id is the owner of the wallet
    amount     numeric(18, 2)           NOT NULL DEFAULT 0 CHECK ( amount >= 0 ), -- amount is the balance of the wallet
    version    bigint                   NOT NULL DEFAULT 0,                       -- version is the version of the wallet
    created_at timestamp with time zone NOT NULL DEFAULT now(),                   -- created_at is the time when the wallet is created
    updated_at timestamp with time zone NOT NULL DEFAULT now()                    -- updated_at is the time when the wallet is updated
);

CREATE INDEX balances_owner_id_idx ON wallet (user_id);

CREATE TYPE transaction_type AS ENUM ('DEPOSIT', 'WITHDRAW', 'PURCHASE', 'REFUND', 'TRANSFER_IN', 'TRANSFER_OUT');
CREATE TYPE initiator_type AS ENUM ('SYSTEM', 'USER', 'BACKOFFICE');
CREATE TYPE reference_type AS ENUM ('ORDER', 'BANK_TRANSACTION', 'TRANSACTION');

CREATE TABLE transactions
(
    id               uuid                      DEFAULT uuid_generate_v4(),
    wallet_id        uuid REFERENCES wallet (id),                   -- wallet_id is the owner of the transaction
    type             transaction_type NOT NULL,                     -- type is the type of the transaction
    amount           numeric(18, 2)   NOT NULL,                     -- amount is the amount of the transaction
    reference_id     text,                                          -- reference_id is the reference of the transaction (e.g. order_id, bank_transaction_id, transaction_id)
    reference_type   reference_type,                                -- reference_type is the type of the reference (e.g. order, bank_transaction, transaction)
    initiated_by     initiator_type   NOT NULL,                     -- initiated_by is the initiator of the transaction
    previous_balance numeric(18, 2)   NOT NULL CHECK (amount >= 0), -- previous_balance is the balance before the transaction
    new_balance      numeric(18, 2)   NOT NULL CHECK (amount >= 0), -- new_balance is the balance after the transaction
    created_at       timestamp with time zone  DEFAULT NOW(),       -- created_at is the time when the transaction is created
    version          bigint           NOT NULL DEFAULT 0,           -- version is the version of the wallet after this transaction
    PRIMARY KEY (id, created_at),                                   -- primary key is the combination of id and created_at
    CHECK ( previous_balance + amount = new_balance ),              -- check if the balance is correct
    CHECK (
        (amount >= 0 AND type IN ('DEPOSIT', 'REFUND', 'TRANSFER_IN')) OR
        (amount <= 0 AND type IN ('WITHDRAW', 'PURCHASE', 'TRANSFER_OUT'))
        ),                                                          -- check if the amount is correct
    CHECK ( type = 'TRANSFER_IN' AND reference_type = 'TRANSACTION' AND reference_id IS NOT NULL),
    CHECK ( type = 'PURCHASE' AND reference_type = 'ORDER' AND reference_id IS NOT NULL),
    CHECK ( type = 'REFUND' AND reference_type = 'ORDER' AND reference_id IS NOT NULL),
    CHECK ( type = 'WITHDRAW' AND reference_type = 'BANK_TRANSACTION' AND reference_id IS NOT NULL),
    CHECK ( type = 'DEPOSIT' AND reference_type = 'BANK_TRANSACTION' AND reference_id IS NULL),
    CHECK ( type = 'TRANSFER_OUT' AND reference_type IS NULL AND reference_id IS NULL)
) PARTITION BY RANGE (created_at);

CREATE INDEX transactions_wallet_id_idx ON transactions (wallet_id);
CREATE INDEX transactions_created_at_idx ON transactions (created_at);

-- Q4 2024: October 1st to December 31st + half of 3rd quarter
CREATE TABLE transactions_2024_q4 PARTITION OF transactions
    FOR VALUES FROM ('2024-07-01 00:00:00+00') TO ('2025-01-01 00:00:00+00');