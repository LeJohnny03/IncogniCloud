CREATE TABLE IF NOT EXISTS server_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    preferences JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);