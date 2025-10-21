-- Migration: Create clicks table
-- Created: 2025-01-27

CREATE TABLE IF NOT EXISTS clicks (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    bannerid INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add foreign key constraint to banners table
ALTER TABLE clicks 
ADD CONSTRAINT fk_clicks_bannerid 
FOREIGN KEY (bannerid) 
REFERENCES banners(id) 
ON DELETE CASCADE 
ON UPDATE CASCADE;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_clicks_bannerid ON clicks(bannerid);
CREATE INDEX IF NOT EXISTS idx_clicks_timestamp ON clicks(timestamp);
CREATE INDEX IF NOT EXISTS idx_clicks_created_at ON clicks(created_at);

-- Create composite index for common queries
CREATE INDEX IF NOT EXISTS idx_clicks_bannerid_timestamp ON clicks(bannerid, timestamp);
