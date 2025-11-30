-- +goose Up
-- Add receipt_url column to donations table for storing PDF receipt URLs
ALTER TABLE donations ADD COLUMN IF NOT EXISTS receipt_url VARCHAR(500);

-- Add comment to explain the column
COMMENT ON COLUMN donations.receipt_url IS 'URL to the PDF receipt stored in S3';

-- +goose Down
-- Remove receipt_url column from donations table
ALTER TABLE donations DROP COLUMN IF EXISTS receipt_url;
