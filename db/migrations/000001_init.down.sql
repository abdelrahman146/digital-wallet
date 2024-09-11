-- Function to delete wallet schemas on deletion
CREATE OR REPLACE FUNCTION drop_wallet_schema(id TEXT) RETURNS VOID AS $$
BEGIN
EXECUTE 'DROP SCHEMA IF EXISTS ' || id || '_wallet CASCADE';
END;
$$ LANGUAGE plpgsql;

-- Drop Wallets table and cascading schemas
DO $$
DECLARE
r RECORD;
BEGIN
    -- Iterate over existing wallet schemas and drop them
FOR r IN (SELECT id FROM wallets) LOOP
        PERFORM drop_wallet_schema(r.id);
END LOOP;
    -- Drop wallets table itself
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS triggers CASCADE;
DROP TABLE IF EXISTS exchange_rates CASCADE;
DROP TABLE IF EXISTS tiers CASCADE;
DROP TABLE IF EXISTS wallets CASCADE;
END $$;

-- Drop account and transaction ID generators
DROP FUNCTION IF EXISTS generate_account_id;
DROP FUNCTION IF EXISTS generate_transaction_id;

DROP FUNCTION IF EXISTS create_wallet_schema;
DROP FUNCTION IF EXISTS drop_wallet_schema;

DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS transaction_actor_type;
