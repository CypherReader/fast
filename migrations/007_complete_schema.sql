-- 007_complete_schema.sql
-- 1. Users Table Updates
ALTER TABLE users
ADD COLUMN IF NOT EXISTS name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS onboarding_completed BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS goal VARCHAR(50),
    ADD COLUMN IF NOT EXISTS fasting_plan VARCHAR(10),
    ADD COLUMN IF NOT EXISTS sex VARCHAR(20),
    ADD COLUMN IF NOT EXISTS height_cm DECIMAL(5, 2),
    ADD COLUMN IF NOT EXISTS current_weight_lbs DECIMAL(5, 2),
    ADD COLUMN IF NOT EXISTS target_weight_lbs DECIMAL(5, 2),
    ADD COLUMN IF NOT EXISTS timezone VARCHAR(50) DEFAULT 'America/New_York',
    ADD COLUMN IF NOT EXISTS units VARCHAR(10) DEFAULT 'imperial',
    ADD COLUMN IF NOT EXISTS stripe_customer_id VARCHAR(100),
    ADD COLUMN IF NOT EXISTS subscription_status VARCHAR(50) DEFAULT 'trial',
    ADD COLUMN IF NOT EXISTS subscription_id VARCHAR(100),
    ADD COLUMN IF NOT EXISTS vault_enabled BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS trial_ends_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS push_notifications_enabled BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS notification_token VARCHAR(255);
CREATE INDEX IF NOT EXISTS idx_users_stripe_customer_id ON users(stripe_customer_id);
-- 2. Subscriptions Table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    plan_type VARCHAR(20) DEFAULT 'core',
    subscription_price DECIMAL(10, 2) DEFAULT 12.00,
    stripe_subscription_id VARCHAR(100),
    status VARCHAR(20),
    current_period_start DATE,
    current_period_end DATE,
    cancel_at_period_end BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- 3. Vault Participations Table
CREATE TABLE IF NOT EXISTS vault_participations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    month_start DATE NOT NULL,
    month_end DATE NOT NULL,
    deposit_amount DECIMAL(10, 2) DEFAULT 20.00,
    fasts_completed INT DEFAULT 0,
    amount_recovered DECIMAL(10, 2) DEFAULT 0.00,
    refund_processed BOOLEAN DEFAULT FALSE,
    refund_date TIMESTAMP WITH TIME ZONE,
    opted_in BOOLEAN DEFAULT TRUE,
    forfeited_amount DECIMAL(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, month_start)
);
CREATE INDEX IF NOT EXISTS idx_vault_participations_user_month ON vault_participations(user_id, month_start);
-- 4. Fasting Sessions Updates (fasts)
ALTER TABLE fasting_sessions
ADD COLUMN IF NOT EXISTS vault_participation_id UUID REFERENCES vault_participations(id),
    ADD COLUMN IF NOT EXISTS planned_duration_hours INT,
    ADD COLUMN IF NOT EXISTS actual_duration_hours DECIMAL(5, 2),
    ADD COLUMN IF NOT EXISTS completed BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS recovery_amount DECIMAL(10, 2) DEFAULT 2.00,
    ADD COLUMN IF NOT EXISTS phase_reached VARCHAR(20),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();
CREATE INDEX IF NOT EXISTS idx_fasting_sessions_start_time ON fasting_sessions(start_time);
-- 5. Friend Networks Table
CREATE TABLE IF NOT EXISTS friend_networks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    friend_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'connected',
    connected_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, friend_id)
);
-- 6. Friend Network Pots Table
CREATE TABLE IF NOT EXISTS friend_network_pots (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES users(id),
    month_start DATE NOT NULL,
    month_end DATE NOT NULL,
    total_pot_amount DECIMAL(10, 2) DEFAULT 0.00,
    winner_id UUID REFERENCES users(id),
    winner_fasts INT DEFAULT 0,
    payout_processed BOOLEAN DEFAULT FALSE,
    payout_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(owner_id, month_start)
);
-- 7. Friend Network Pot Participants Table
CREATE TABLE IF NOT EXISTS friend_network_pot_participants (
    id UUID PRIMARY KEY,
    pot_id UUID NOT NULL REFERENCES friend_network_pots(id),
    user_id UUID NOT NULL REFERENCES users(id),
    fasts_completed INT DEFAULT 0,
    rank INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(pot_id, user_id)
);
-- 8. Friend Challenges Table
CREATE TABLE IF NOT EXISTS friend_challenges (
    id UUID PRIMARY KEY,
    challenger_id UUID NOT NULL REFERENCES users(id),
    challenged_id UUID NOT NULL REFERENCES users(id),
    month_start DATE NOT NULL,
    month_end DATE NOT NULL,
    pot_amount DECIMAL(10, 2) DEFAULT 40.00,
    challenger_deposit DECIMAL(10, 2) DEFAULT 20.00,
    challenged_deposit DECIMAL(10, 2) DEFAULT 20.00,
    challenger_fasts INT DEFAULT 0,
    challenged_fasts INT DEFAULT 0,
    winner_id UUID REFERENCES users(id),
    status VARCHAR(20),
    payout_amount DECIMAL(10, 2),
    payout_processed BOOLEAN DEFAULT FALSE,
    payout_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_friend_challenges_status ON friend_challenges(status);
-- 9. Tribes Table
CREATE TABLE IF NOT EXISTS tribes (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users(id),
    member_count INT DEFAULT 0,
    is_private BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- 10. Tribe Members Table
CREATE TABLE IF NOT EXISTS tribe_members (
    id UUID PRIMARY KEY,
    tribe_id UUID NOT NULL REFERENCES tribes(id),
    user_id UUID NOT NULL REFERENCES users(id),
    is_admin BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tribe_id, user_id)
);
-- 11. Tribe Pools Table
CREATE TABLE IF NOT EXISTS tribe_pools (
    id UUID PRIMARY KEY,
    tribe_id UUID NOT NULL REFERENCES tribes(id),
    month_start DATE NOT NULL,
    month_end DATE NOT NULL,
    total_pot DECIMAL(10, 2) DEFAULT 0.00,
    participant_count INT DEFAULT 0,
    status VARCHAR(20),
    first_place_payout DECIMAL(10, 2),
    second_place_payout DECIMAL(10, 2),
    third_place_payout DECIMAL(10, 2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tribe_id, month_start)
);
-- 12. Tribe Pool Participants Table
CREATE TABLE IF NOT EXISTS tribe_pool_participants (
    id UUID PRIMARY KEY,
    pool_id UUID NOT NULL REFERENCES tribe_pools(id),
    user_id UUID NOT NULL REFERENCES users(id),
    deposit_amount DECIMAL(10, 2) DEFAULT 20.00,
    fasts_completed INT DEFAULT 0,
    final_rank INT,
    payout_amount DECIMAL(10, 2) DEFAULT 0.00,
    payout_processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(pool_id, user_id)
);
-- 13. Weight Logs Table
CREATE TABLE IF NOT EXISTS weight_logs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    weight_lbs DECIMAL(5, 2) NOT NULL,
    weight_kg DECIMAL(5, 2) NOT NULL,
    logged_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- 14. Meals Table Updates
ALTER TABLE meals
ADD COLUMN IF NOT EXISTS meal_type VARCHAR(20),
    ADD COLUMN IF NOT EXISTS is_keto BOOLEAN;
-- 15. Hydration Logs Table
CREATE TABLE IF NOT EXISTS hydration_logs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    glasses_count INT DEFAULT 0,
    logged_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, logged_date)
);
-- 16. Activity Logs Table
CREATE TABLE IF NOT EXISTS activity_logs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    steps INT DEFAULT 0,
    distance_km DECIMAL(5, 2),
    calories_burned INT,
    logged_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, logged_date)
);
-- 17. Notifications Table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(50),
    title VARCHAR(200),
    message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Add columns that might not exist
ALTER TABLE notifications
ADD COLUMN IF NOT EXISTS read BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS link VARCHAR(500);
-- CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, read);
-- 18. Commitment Contracts Table
CREATE TABLE IF NOT EXISTS commitment_contracts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    commitment_text TEXT,
    goals JSONB,
    signed_at TIMESTAMP WITH TIME ZONE
);