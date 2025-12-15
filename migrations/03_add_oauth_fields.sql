-- Add OAuth fields to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS google_id VARCHAR(255) UNIQUE,
    ADD COLUMN IF NOT EXISTS oauth_provider VARCHAR(50),
    ADD COLUMN IF NOT EXISTS profile_picture_url TEXT;
-- Make password_hash nullable for OAuth-only users
ALTER TABLE users
ALTER COLUMN password_hash DROP NOT NULL;
-- Create index for faster Google ID lookups
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
-- Create index for OAuth provider
CREATE INDEX IF NOT EXISTS idx_users_oauth_provider ON users(oauth_provider);