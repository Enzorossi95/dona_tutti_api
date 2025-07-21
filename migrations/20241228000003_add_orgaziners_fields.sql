-- +goose Up
-- Add organizers fields
ALTER TABLE organizers 
ADD COLUMN email VARCHAR(255) NULL,
ADD COLUMN phone INTEGER NULL CHECK (phone IS NULL OR phone >= 0),
ADD COLUMN website VARCHAR(255) NULL,
ADD COLUMN address VARCHAR(500) NULL;

-- +goose Down

-- Remove the added columns
ALTER TABLE organizers 
DROP COLUMN IF EXISTS email,
DROP COLUMN IF EXISTS phone,
DROP COLUMN IF EXISTS website,
DROP COLUMN IF EXISTS address; 