-- Migration: 010_tribes
-- Description: Create tribes and tribe_memberships tables for social groups feature
-- Author: System
-- Date: 2025-12-16
-- Create tribes table
CREATE TABLE IF NOT EXISTS tribes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(60) UNIQUE NOT NULL,
    description VARCHAR(500) NOT NULL,
    avatar_url TEXT,
    cover_photo_url TEXT,
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fasting_schedule VARCHAR(20) NOT NULL,
    primary_goal VARCHAR(30) NOT NULL,
    category JSONB DEFAULT '[]'::jsonb,
    privacy VARCHAR(20) NOT NULL DEFAULT 'public',
    rules TEXT,
    member_count INTEGER NOT NULL DEFAULT 0,
    active_member_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS idx_tribes_slug ON tribes(slug);
CREATE INDEX IF NOT EXISTS idx_tribes_creator ON tribes(creator_id);
CREATE INDEX IF NOT EXISTS idx_tribes_privacy ON tribes(privacy);
CREATE INDEX IF NOT EXISTS idx_tribes_schedule ON tribes(fasting_schedule);
CREATE INDEX IF NOT EXISTS idx_tribes_created ON tribes(created_at DESC);
-- Create tribe_memberships table
CREATE TABLE IF NOT EXISTS tribe_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tribe_id UUID NOT NULL REFERENCES tribes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    left_at TIMESTAMP WITH TIME ZONE,
    notifications_enabled BOOLEAN DEFAULT TRUE,
    UNIQUE(tribe_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_memberships_user ON tribe_memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_memberships_tribe ON tribe_memberships(tribe_id);
CREATE INDEX IF NOT EXISTS idx_memberships_status ON tribe_memberships(status);
CREATE INDEX IF NOT EXISTS idx_memberships_joined ON tribe_memberships(joined_at DESC);