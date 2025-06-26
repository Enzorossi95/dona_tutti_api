-- +goose Up
-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    resource VARCHAR(50) NOT NULL, -- campaigns, donations, users, etc.
    action VARCHAR(20) NOT NULL,   -- create, read, update, delete
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create role_permissions junction table (many-to-many)
CREATE TABLE IF NOT EXISTS role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);

-- Add role_id to users table
ALTER TABLE users ADD COLUMN role_id UUID REFERENCES roles(id);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_roles_active ON roles(is_active);
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource, action);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);

-- Insert default roles
INSERT INTO roles (id, name, description) VALUES 
    ('11111111-1111-1111-1111-111111111111', 'admin', 'Administrator with full system access'),
    ('22222222-2222-2222-2222-222222222222', 'donor', 'Donor who can create and manage donations'),
    ('33333333-3333-3333-3333-333333333333', 'guest', 'Guest with read-only access to public content');

-- Insert default permissions
INSERT INTO permissions (name, resource, action, description) VALUES 
    -- Campaign permissions
    ('campaigns:create', 'campaigns', 'create', 'Create new campaigns'),
    ('campaigns:read', 'campaigns', 'read', 'Read campaign information'),
    ('campaigns:update', 'campaigns', 'update', 'Update existing campaigns'),
    ('campaigns:delete', 'campaigns', 'delete', 'Delete campaigns'),
    
    -- Donation permissions
    ('donations:create', 'donations', 'create', 'Create new donations'),
    ('donations:read', 'donations', 'read', 'Read donation information'),
    ('donations:update', 'donations', 'update', 'Update existing donations'),
    ('donations:delete', 'donations', 'delete', 'Delete donations'),
    
    -- User permissions
    ('users:create', 'users', 'create', 'Create new users'),
    ('users:read', 'users', 'read', 'Read user information'),
    ('users:update', 'users', 'update', 'Update user information'),
    ('users:delete', 'users', 'delete', 'Delete users'),
    
    -- Category permissions
    ('categories:create', 'categories', 'create', 'Create new categories'),
    ('categories:read', 'categories', 'read', 'Read category information'),
    ('categories:update', 'categories', 'update', 'Update existing categories'),
    ('categories:delete', 'categories', 'delete', 'Delete categories'),
    
    -- Organizer permissions
    ('organizers:create', 'organizers', 'create', 'Create new organizers'),
    ('organizers:read', 'organizers', 'read', 'Read organizer information'),
    ('organizers:update', 'organizers', 'update', 'Update existing organizers'),
    ('organizers:delete', 'organizers', 'delete', 'Delete organizers'),
    
    -- Donor permissions
    ('donors:create', 'donors', 'create', 'Create new donor profiles'),
    ('donors:read', 'donors', 'read', 'Read donor information'),
    ('donors:update', 'donors', 'update', 'Update donor profiles'),
    ('donors:delete', 'donors', 'delete', 'Delete donor profiles');

-- Assign permissions to roles
-- Admin gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT '11111111-1111-1111-1111-111111111111', id FROM permissions;

-- Donor permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT '22222222-2222-2222-2222-222222222222', id FROM permissions 
WHERE name IN (
    'campaigns:read', 
    'categories:read', 
    'organizers:read',
    'donations:create', 
    'donations:read', 
    'donations:update', 
    'donations:delete',
    'donors:create',
    'donors:read',
    'donors:update',
    'users:update'
);

-- Guest permissions (read-only public content)
INSERT INTO role_permissions (role_id, permission_id)
SELECT '33333333-3333-3333-3333-333333333333', id FROM permissions 
WHERE name IN (
    'campaigns:read', 
    'categories:read', 
    'organizers:read'
);

-- Set default role for existing users (guest)
UPDATE users SET role_id = '33333333-3333-3333-3333-333333333333' WHERE role_id IS NULL;

-- Make role_id NOT NULL after setting defaults
ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;

-- +goose Down
-- Remove role_id column from users
ALTER TABLE users DROP COLUMN role_id;

-- Drop tables in reverse order
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;