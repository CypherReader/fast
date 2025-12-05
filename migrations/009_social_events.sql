-- 009_social_events.sql
CREATE TABLE IF NOT EXISTS social_events (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    event_type VARCHAR(50) NOT NULL,
    data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    likes INT DEFAULT 0,
    comments INT DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_social_events_created_at ON social_events(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_social_events_user_id ON social_events(user_id);