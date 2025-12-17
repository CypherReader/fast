-- SOS Flares table
CREATE TABLE IF NOT EXISTS sos_flares (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fasting_id UUID REFERENCES fasting_sessions(id) ON DELETE
    SET NULL,
        tribe_id UUID REFERENCES tribes(id) ON DELETE
    SET NULL,
        description TEXT,
        hours_fasted DECIMAL(10, 2),
        status VARCHAR(20) DEFAULT 'active',
        hype_count INT DEFAULT 0,
        is_anonymous BOOLEAN DEFAULT false,
        cortex_responded BOOLEAN DEFAULT false,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        resolved_at TIMESTAMP,
        INDEX idx_sos_user (user_id),
        INDEX idx_sos_status (status),
        INDEX idx_sos_created (created_at DESC),
        INDEX idx_sos_active (user_id, status) -- For finding active SOS
);
-- Hype responses table
CREATE TABLE IF NOT EXISTS hype_responses (
    id UUID PRIMARY KEY,
    sos_id UUID NOT NULL REFERENCES sos_flares(id) ON DELETE CASCADE,
    from_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_name VARCHAR(255),
    message TEXT,
    emoji VARCHAR(10),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_hype_sos (sos_id),
    INDEX idx_hype_from (from_user_id),
    INDEX idx_hype_daily (from_user_id, created_at) -- For daily limits
);
-- User SOS settings (add to users table)
ALTER TABLE users
ADD COLUMN IF NOT EXISTS notify_tribe_on_sos BOOLEAN DEFAULT true;
ALTER TABLE users
ADD COLUMN IF NOT EXISTS sos_anonymous_mode BOOLEAN DEFAULT false;
ALTER TABLE users
ADD COLUMN IF NOT EXISTS last_sos_at TIMESTAMP;
-- Create index for rate limiting
CREATE INDEX IF NOT EXISTS idx_users_last_sos ON users(last_sos_at);
-- Notification types (update notification_types enum if using PostgreSQL enums)
-- For now, we'll handle this in the application layer with string constants