-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_account_id VARCHAR(255) NOT NULL,
    encrypted_token BYTEA NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(owner_user_id, provider, provider_account_id)
);

-- Create index on owner_user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_accounts_owner_user_id ON accounts(owner_user_id);

-- Create index on provider for filtering
CREATE INDEX IF NOT EXISTS idx_accounts_provider ON accounts(provider);

-- Create composite index for provider-specific lookups
CREATE INDEX IF NOT EXISTS idx_accounts_owner_provider ON accounts(owner_user_id, provider);
