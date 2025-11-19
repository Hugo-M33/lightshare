-- Drop indexes
DROP INDEX IF EXISTS idx_users_stripe_customer_id;
DROP INDEX IF EXISTS idx_users_email;

-- Drop users table
DROP TABLE IF EXISTS users;
