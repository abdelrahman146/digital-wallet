DROP TABLE IF EXISTS transactions_2024_q4;

DROP INDEX IF EXISTS transactions_created_at_idx;
DROP INDEX IF EXISTS transactions_wallet_id_idx;

DROP TABLE IF EXISTS transactions;

DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS initiator_type;
DROP TYPE IF EXISTS reference_type;

DROP INDEX IF EXISTS balances_owner_id_idx;

DROP TABLE IF EXISTS wallet;
