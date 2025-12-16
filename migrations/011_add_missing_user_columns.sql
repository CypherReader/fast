-- Migration: 011_add_missing_user_columns
-- Description: Add missing user columns that are referenced in code but not in schema
-- Author: System
-- Date: 2025-12-16
ALTER TABLE users
ADD COLUMN IF NOT EXISTS vault_deposit DECIMAL(10, 2) DEFAULT 20.00,
    ADD COLUMN IF NOT EXISTS earned_refund DECIMAL(10, 2) DEFAULT 0.00,
    ADD COLUMN IF NOT EXISTS tribe_id UUID REFERENCES tribes(id);
CREATE INDEX IF NOT EXISTS idx_users_tribe_id ON users(tribe_id);