-- +goose Up
-- Create donors table
CREATE TABLE IF NOT EXISTS donors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT false,
    phone VARCHAR(50),
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_donors_email ON donors(email);

-- Create donation types and table
DO $donation_status_type$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'donation_status') THEN CREATE TYPE donation_status AS ENUM ('completed', 'pending', 'failed', 'refunded'); END IF; END $donation_status_type$;

CREATE TABLE IF NOT EXISTS donations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id),
    donor_id UUID NOT NULL REFERENCES donors(id),
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    message TEXT,
    is_anonymous BOOLEAN DEFAULT false,
    payment_method VARCHAR(255) NOT NULL,
    status donation_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_donations_campaign_id ON donations(campaign_id);
CREATE INDEX IF NOT EXISTS idx_donations_donor_id ON donations(donor_id);
CREATE INDEX IF NOT EXISTS idx_donations_date ON donations(date);
CREATE INDEX IF NOT EXISTS idx_donations_status ON donations(status);

-- +goose Down
DROP TABLE IF EXISTS donations;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS donation_status;
DROP TABLE IF EXISTS donors; 