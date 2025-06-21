-- +goose Up
-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create articles table
CREATE TABLE IF NOT EXISTS articles (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create campaign_categories table
CREATE TABLE IF NOT EXISTS campaign_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create organizers table
CREATE TABLE IF NOT EXISTS organizers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    avatar VARCHAR(500),
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create campaigns table
CREATE TABLE IF NOT EXISTS campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image VARCHAR(500),
    goal DECIMAL(12,2) NOT NULL CHECK (goal > 0),
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    location VARCHAR(255),
    category_id UUID REFERENCES campaign_categories(id),
    urgency INTEGER CHECK (urgency >= 1 AND urgency <= 10),
    organizer_id UUID REFERENCES organizers(id),
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_end_date_after_start CHECK (end_date > start_date)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_campaigns_category_id ON campaigns(category_id);
CREATE INDEX IF NOT EXISTS idx_campaigns_organizer_id ON campaigns(organizer_id);
CREATE INDEX IF NOT EXISTS idx_campaigns_status ON campaigns(status);
CREATE INDEX IF NOT EXISTS idx_campaigns_start_date ON campaigns(start_date);
CREATE INDEX IF NOT EXISTS idx_campaigns_end_date ON campaigns(end_date);

-- Insert sample data for campaign categories
INSERT INTO campaign_categories (id, name, description) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Education', 'Campaigns related to education and learning'),
    ('550e8400-e29b-41d4-a716-446655440002', 'Health', 'Medical and health-related campaigns'),
    ('550e8400-e29b-41d4-a716-446655440003', 'Environment', 'Environmental protection and sustainability'),
    ('550e8400-e29b-41d4-a716-446655440004', 'Community', 'Community development and social causes')
ON CONFLICT (name) DO NOTHING;

-- Insert sample data for organizers
INSERT INTO organizers (id, name, avatar, verified) VALUES
    ('660e8400-e29b-41d4-a716-446655440001', 'Education Foundation', 'https://example.com/avatars/education-foundation.jpg', true),
    ('660e8400-e29b-41d4-a716-446655440002', 'Medical Relief Org', 'https://example.com/avatars/medical-relief.jpg', true),
    ('660e8400-e29b-41d4-a716-446655440003', 'Water for All', 'https://example.com/avatars/water-for-all.jpg', false)
ON CONFLICT (id) DO NOTHING;

-- Insert sample articles
INSERT INTO articles (id, title, content) VALUES
    ('1', 'Primer artículo', 'Contenido del primer artículo'),
    ('2', 'Segundo artículo', 'Contenido del segundo artículo'),
    ('3', 'Tercer artículo', 'Contenido del tercer artículo')
ON CONFLICT (id) DO NOTHING;

-- Insert sample campaigns
INSERT INTO campaigns (id, title, description, image, goal, start_date, end_date, location, category_id, urgency, organizer_id, status) VALUES
    (
        '770e8400-e29b-41d4-a716-446655440001',
        'Help Build School in Rural Area',
        'We need funds to build a new school for children in remote villages',
        'https://example.com/school.jpg',
        50000.00,
        CURRENT_TIMESTAMP - INTERVAL '10 days',
        CURRENT_TIMESTAMP + INTERVAL '30 days',
        'Rural Village, State',
        '550e8400-e29b-41d4-a716-446655440001',
        8,
        '660e8400-e29b-41d4-a716-446655440001',
        'active'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440002',
        'Emergency Medical Fund',
        'Urgent medical treatment needed for local community',
        'https://example.com/medical.jpg',
        25000.00,
        CURRENT_TIMESTAMP - INTERVAL '5 days',
        CURRENT_TIMESTAMP + INTERVAL '25 days',
        'City Hospital',
        '550e8400-e29b-41d4-a716-446655440002',
        10,
        '660e8400-e29b-41d4-a716-446655440002',
        'active'
    ),
    (
        '770e8400-e29b-41d4-a716-446655440003',
        'Clean Water Initiative',
        'Providing clean water access to underserved communities',
        'https://example.com/water.jpg',
        75000.00,
        CURRENT_TIMESTAMP - INTERVAL '15 days',
        CURRENT_TIMESTAMP + INTERVAL '60 days',
        'Multiple Locations',
        '550e8400-e29b-41d4-a716-446655440003',
        7,
        '660e8400-e29b-41d4-a716-446655440003',
        'active'
    )
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS campaigns;
DROP TABLE IF EXISTS organizers;
DROP TABLE IF EXISTS campaign_categories;
DROP TABLE IF EXISTS articles;
