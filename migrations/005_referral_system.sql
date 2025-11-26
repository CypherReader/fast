-- Add referral_code to users table
ALTER TABLE users
ADD COLUMN referral_code VARCHAR(10) UNIQUE;
-- Create referrals table
CREATE TABLE referrals (
    id UUID PRIMARY KEY,
    referrer_id UUID NOT NULL REFERENCES users(id),
    referee_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL,
    -- pending, completed
    reward_value DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(referee_id) -- A user can only be referred once
);
CREATE INDEX idx_referrals_referrer_id ON referrals(referrer_id);