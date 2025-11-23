-- Drop indexes
DROP INDEX IF EXISTS idx_accounts_owner_provider;
DROP INDEX IF EXISTS idx_accounts_provider;
DROP INDEX IF EXISTS idx_accounts_owner_user_id;

-- Drop accounts table
DROP TABLE IF EXISTS accounts;
