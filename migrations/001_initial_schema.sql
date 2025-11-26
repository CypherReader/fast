-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    subscription_tier VARCHAR(50) NOT NULL DEFAULT 'Free',
    discipline_index FLOAT NOT NULL DEFAULT 0,
    current_price FLOAT NOT NULL DEFAULT 50.0,
    signed_contract BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS fasting_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    goal_hours FLOAT NOT NULL,
    plan_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE INDEX idx_fasting_sessions_user_id ON fasting_sessions(user_id);
CREATE TABLE IF NOT EXISTS keto_entries (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    logged_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ketone_level FLOAT NOT NULL,
    acetone_level FLOAT NOT NULL,
    source VARCHAR(50) NOT NULL,
    CONSTRAINT fk_keto_user FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE INDEX idx_keto_entries_user_id ON keto_entries(user_id);
CREATE TABLE IF NOT EXISTS meals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    calories INTEGER,
    protein FLOAT,
    carbs FLOAT,
    fats FLOAT,
    image_url VARCHAR(255),
    logged_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT fk_meal_user FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE INDEX idx_meals_user_id ON meals(user_id);
CREATE TABLE IF NOT EXISTS telemetry_data (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    metric_type VARCHAR(50) NOT NULL,
    value FLOAT NOT NULL,
    source VARCHAR(50),
    CONSTRAINT fk_telemetry_user FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE INDEX idx_telemetry_user_id ON telemetry_data(user_id);