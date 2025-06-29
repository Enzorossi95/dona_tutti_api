-- +goose Up
-- Create payment methods system

-- 1. Create payment_methods table (catalog of available payment methods)
CREATE TABLE IF NOT EXISTS payment_methods (
    id SERIAL PRIMARY KEY,
    code VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create campaign_payment_methods table (N-N relationship campaigns <-> payment methods)
CREATE TABLE IF NOT EXISTS campaign_payment_methods (
    id SERIAL PRIMARY KEY,
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    payment_method_id INTEGER NOT NULL REFERENCES payment_methods(id) ON DELETE CASCADE,
    instructions TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(campaign_id, payment_method_id)
);

-- 3. Create transfer_details table (specific details for bank transfers)
CREATE TABLE IF NOT EXISTS transfer_details (
    id SERIAL PRIMARY KEY,
    campaign_payment_method_id INTEGER NOT NULL
        REFERENCES campaign_payment_methods(id)
        ON DELETE CASCADE,
    bank_name VARCHAR(100) NOT NULL,
    account_holder VARCHAR(100) NOT NULL,
    cbu VARCHAR(22) NOT NULL,
    alias VARCHAR(30),
    swift_code VARCHAR(11),
    additional_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. Create cash_locations table (specific details for cash payments)
CREATE TABLE IF NOT EXISTS cash_locations (
    id SERIAL PRIMARY KEY,
    campaign_payment_method_id INTEGER NOT NULL
        REFERENCES campaign_payment_methods(id)
        ON DELETE CASCADE,
    location_name VARCHAR(100) NOT NULL,
    address VARCHAR(200) NOT NULL,
    contact_info VARCHAR(100),
    available_hours VARCHAR(100),
    additional_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_campaign_payment_methods_campaign_id ON campaign_payment_methods(campaign_id);
CREATE INDEX IF NOT EXISTS idx_campaign_payment_methods_payment_method_id ON campaign_payment_methods(payment_method_id);
CREATE INDEX IF NOT EXISTS idx_transfer_details_campaign_payment_method_id ON transfer_details(campaign_payment_method_id);
CREATE INDEX IF NOT EXISTS idx_cash_locations_campaign_payment_method_id ON cash_locations(campaign_payment_method_id);

-- Insert initial payment methods data
INSERT INTO payment_methods (code, name, is_active) VALUES
    ('transfer', 'Transferencia Bancaria', true),
    ('cash', 'Efectivo', true),
    ('mercadopago', 'MercadoPago', true)
ON CONFLICT (code) DO NOTHING;

-- +goose Down
-- Drop tables in reverse order due to foreign key constraints
DROP TABLE IF EXISTS cash_locations;
DROP TABLE IF EXISTS transfer_details;
DROP TABLE IF EXISTS campaign_payment_methods;
DROP TABLE IF EXISTS payment_methods;