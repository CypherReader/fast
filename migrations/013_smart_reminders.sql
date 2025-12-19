-- Smart Reminders tables and user preferences
-- Scheduled reminders table
CREATE TABLE IF NOT EXISTS scheduled_reminders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reminder_type VARCHAR(50) NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    sent BOOLEAN DEFAULT false,
    message TEXT,
    data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_reminders_user ON scheduled_reminders(user_id);
CREATE INDEX IF NOT EXISTS idx_reminders_scheduled ON scheduled_reminders(scheduled_at, sent);
CREATE INDEX IF NOT EXISTS idx_reminders_pending ON scheduled_reminders(sent, scheduled_at)
WHERE sent = false;
-- Add reminder preferences to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS reminder_fast_start BOOLEAN DEFAULT true;
ALTER TABLE users
ADD COLUMN IF NOT EXISTS reminder_fast_end BOOLEAN DEFAULT true;
ALTER TABLE users
ADD COLUMN IF NOT EXISTS reminder_hydration BOOLEAN DEFAULT false;
ALTER TABLE users
ADD COLUMN IF NOT EXISTS preferred_fast_start_hour INT DEFAULT 20;
-- 8 PM default
ALTER TABLE users
ADD COLUMN IF NOT EXISTS hydration_interval_minutes INT DEFAULT 60;