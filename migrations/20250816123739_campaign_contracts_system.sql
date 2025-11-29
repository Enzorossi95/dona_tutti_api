-- +goose Up
-- +goose StatementBegin
-- Create campaign status type
DO $campaign_status_type$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'campaign_status') THEN 
        CREATE TYPE campaign_status AS ENUM (
            'draft',
            'pending_approval', 
            'active',
            'paused',
            'completed',
            'rejected'
        );
    END IF; 
END $campaign_status_type$;
-- +goose StatementEnd

-- Alter campaigns table to use the new enum type
-- First, update existing data to use valid enum values
UPDATE campaigns SET status = 'active' WHERE status IS NULL OR status = '';

-- Add a temporary column with the new enum type
ALTER TABLE campaigns ADD COLUMN status_new campaign_status;

-- Copy data from old column to new column (with default value for NULL)
UPDATE campaigns SET status_new = 
    CASE 
        WHEN status = 'draft' THEN 'draft'::campaign_status
        WHEN status = 'pending_approval' THEN 'pending_approval'::campaign_status
        WHEN status = 'paused' THEN 'paused'::campaign_status
        WHEN status = 'completed' THEN 'completed'::campaign_status
        WHEN status = 'rejected' THEN 'rejected'::campaign_status
        ELSE 'active'::campaign_status
    END;

-- Drop the old column and rename the new one
ALTER TABLE campaigns DROP COLUMN status;
ALTER TABLE campaigns RENAME COLUMN status_new TO status;

-- Set default and not null constraint
ALTER TABLE campaigns ALTER COLUMN status SET DEFAULT 'draft'::campaign_status;
ALTER TABLE campaigns ALTER COLUMN status SET NOT NULL;

-- Create campaign_contracts table
CREATE TABLE IF NOT EXISTS campaign_contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    organizer_id UUID NOT NULL REFERENCES organizers(id),
    contract_pdf_url TEXT NOT NULL,
    contract_hash VARCHAR(64) NOT NULL,
    accepted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    acceptance_ip VARCHAR(45) NOT NULL,
    acceptance_user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_campaign_contract UNIQUE(campaign_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_campaign_contracts_campaign_id ON campaign_contracts(campaign_id);
CREATE INDEX IF NOT EXISTS idx_campaign_contracts_organizer_id ON campaign_contracts(organizer_id);
CREATE INDEX IF NOT EXISTS idx_campaign_contracts_accepted_at ON campaign_contracts(accepted_at);

-- +goose Down
-- Drop the campaign_contracts table
DROP TABLE IF EXISTS campaign_contracts;

-- Revert campaigns.status to VARCHAR
ALTER TABLE campaigns ADD COLUMN status_old VARCHAR(50);
UPDATE campaigns SET status_old = status::TEXT;
ALTER TABLE campaigns DROP COLUMN status;
ALTER TABLE campaigns RENAME COLUMN status_old TO status;
ALTER TABLE campaigns ALTER COLUMN status SET DEFAULT 'active';

-- Drop the enum type
DROP TYPE IF EXISTS campaign_status;

