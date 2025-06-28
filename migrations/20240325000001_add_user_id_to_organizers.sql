-- +goose Up
-- Add user_id column to organizers table
ALTER TABLE organizers 
ADD COLUMN user_id UUID NULL;

-- Add foreign key constraint to users table
ALTER TABLE organizers 
ADD CONSTRAINT fk_organizers_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Create index on user_id for performance
CREATE INDEX IF NOT EXISTS idx_organizers_user_id ON organizers(user_id);

-- +goose Down
-- Remove the index
DROP INDEX IF EXISTS idx_organizers_user_id;

-- Remove the foreign key constraint
ALTER TABLE organizers 
DROP CONSTRAINT IF EXISTS fk_organizers_user_id;

-- Remove the user_id column
ALTER TABLE organizers 
DROP COLUMN IF EXISTS user_id;