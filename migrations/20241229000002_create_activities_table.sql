-- +goose Up
-- Create activities table
CREATE TABLE IF NOT EXISTS activities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    date TIMESTAMP NOT NULL,
    type VARCHAR(100) NOT NULL,
    author VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index for better performance when querying by campaign_id
CREATE INDEX IF NOT EXISTS idx_activities_campaign_id ON activities(campaign_id);

-- Create index for querying by date
CREATE INDEX IF NOT EXISTS idx_activities_date ON activities(date);

-- +goose Down
-- Drop the activities table
DROP INDEX IF EXISTS idx_activities_date;
DROP INDEX IF EXISTS idx_activities_campaign_id;
DROP TABLE IF EXISTS activities;