-- +goose Up
-- Fix organizer phone field type from INTEGER to VARCHAR
-- First, drop the CHECK constraint that prevents the type change
ALTER TABLE organizers 
DROP CONSTRAINT IF EXISTS organizers_phone_check;

-- Now change the column type
ALTER TABLE organizers 
ALTER COLUMN phone TYPE VARCHAR(20);

-- +goose Down
-- Revert phone field back to INTEGER (note: this may cause data loss if phone numbers exceed INTEGER range)
ALTER TABLE organizers 
ALTER COLUMN phone TYPE INTEGER USING phone::INTEGER;

-- Re-add the CHECK constraint
ALTER TABLE organizers 
ADD CONSTRAINT organizers_phone_check CHECK (phone IS NULL OR phone >= 0);