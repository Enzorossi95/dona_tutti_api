-- +goose Up
-- Add beneficiary and urgency context fields to campaigns table (all nullable)
ALTER TABLE campaigns 
ADD COLUMN beneficiary_name VARCHAR(255) NULL,
ADD COLUMN beneficiary_age INTEGER NULL CHECK (beneficiary_age IS NULL OR beneficiary_age >= 0),
ADD COLUMN current_situation TEXT NULL,
ADD COLUMN urgency_reason VARCHAR(500) NULL;

-- Create indexes for potential queries on these fields
CREATE INDEX IF NOT EXISTS idx_campaigns_beneficiary_name ON campaigns(beneficiary_name);
CREATE INDEX IF NOT EXISTS idx_campaigns_beneficiary_age ON campaigns(beneficiary_age);

-- +goose Down
-- Remove indexes first
DROP INDEX IF EXISTS idx_campaigns_beneficiary_age;
DROP INDEX IF EXISTS idx_campaigns_beneficiary_name;

-- Remove the added columns
ALTER TABLE campaigns 
DROP COLUMN IF EXISTS urgency_reason,
DROP COLUMN IF EXISTS current_situation,
DROP COLUMN IF EXISTS beneficiary_age,
DROP COLUMN IF EXISTS beneficiary_name; 