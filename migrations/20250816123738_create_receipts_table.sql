-- +goose Up
-- Create receipts table for campaign receipts management
CREATE TABLE IF NOT EXISTS receipts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    provider VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total DECIMAL(12,2) NOT NULL CHECK (total > 0),
    quantity INTEGER DEFAULT 1 CHECK (quantity >= 1),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    document_url VARCHAR(500),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_receipts_campaign_id ON receipts(campaign_id);
CREATE INDEX IF NOT EXISTS idx_receipts_date ON receipts(date);
CREATE INDEX IF NOT EXISTS idx_receipts_provider ON receipts(provider);

-- Add comment to table
COMMENT ON TABLE receipts IS 'Stores receipt documents for campaign expenses and transactions';
COMMENT ON COLUMN receipts.provider IS 'Name of the provider or vendor';
COMMENT ON COLUMN receipts.total IS 'Total amount of the receipt';
COMMENT ON COLUMN receipts.quantity IS 'Number of items or units';
COMMENT ON COLUMN receipts.document_url IS 'URL to the PDF document stored in S3';
COMMENT ON COLUMN receipts.note IS 'Additional notes or comments about the receipt';

-- +goose Down
DROP TABLE IF EXISTS receipts;