# Dona Tutti API - Backend Application Information

## Table of Contents
1. [Overview](#overview)
2. [API Endpoints & Examples](#api-endpoints--examples)
3. [Database Schemas](#database-schemas)
4. [Authentication Flow](#authentication-flow)
5. [Business Logic](#business-logic)

---

## Overview

**Dona Tutti API** is a comprehensive donation platform backend built with:
- **Language**: Go 1.24
- **Framework**: Echo v4 (HTTP web framework)
- **ORM**: GORM with PostgreSQL
- **Database**: PostgreSQL 15
- **Authentication**: JWT (HS256)
- **Authorization**: RBAC (Role-Based Access Control)
- **File Storage**: AWS S3
- **Container**: Docker & Docker Compose

### Architecture Pattern
Domain-Driven Design with clean architecture:
```
/{domain}/
├── model.go        # Database models with GORM tags
├── repository.go   # Data access layer
├── service.go      # Business logic layer
├── handlers.go     # HTTP handlers (Echo)
└── {domain}.go     # Domain entities and validation
```

---

## API Endpoints & Examples

### Base URL
```
http://localhost:9999/api
```

### Authentication Endpoints

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "strongpassword123",
  "first_name": "John",
  "last_name": "Doe"
}

Response: 201 Created
{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "strongpassword123"
}

Response: 200 OK
{
  "token": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

#### Request Password Reset
```http
POST /api/auth/password-reset/request
Content-Type: application/json

{
  "email": "user@example.com"
}

Response: 200 OK
{
  "message": "Password reset instructions sent to email"
}
```

#### Reset Password
```http
POST /api/auth/password-reset/reset
Content-Type: application/json

{
  "token": "reset-token-123",
  "new_password": "newpassword123"
}

Response: 200 OK
{
  "message": "Password reset successfully"
}
```

#### Check Permission
```http
POST /api/auth/check-permission
Authorization: Bearer {token}
Content-Type: application/json

{
  "permission": "campaigns:create"
}

Response: 200 OK
{
  "has_permission": true
}
```

### User Endpoints

#### Get Current User
```http
GET /api/users/me
Authorization: Bearer {token}

Response: 200 OK
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@email.com",
  "first_name": "Juan",
  "last_name": "Pérez",
  "role": {
    "id": "11111111-1111-1111-1111-111111111111",
    "name": "admin"
  },
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### List Users (Admin Only)
```http
GET /api/users
Authorization: Bearer {token}

Response: 200 OK
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role_id": "22222222-2222-2222-2222-222222222222",
    "is_active": true
  }
]
```

#### Update User
```http
PUT /api/users/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith"
}

Response: 200 OK
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "first_name": "Jane",
    "last_name": "Smith"
  }
}
```

#### Update Password
```http
PUT /api/users/{id}/password
Authorization: Bearer {token}
Content-Type: application/json

{
  "current_password": "oldpassword123",
  "new_password": "newpassword123"
}

Response: 200 OK
{
  "message": "Password updated successfully"
}
```

### Campaign Endpoints

#### Get Campaign (Public)
```http
GET /api/campaigns/{id}

Response: 200 OK
{
  "id": "770e8400-e29b-41d4-a716-446655440001",
  "title": "Help Build School in Rural Area",
  "description": "We need funds to build a new school",
  "image": "https://example.com/school.jpg",
  "goal": 50000.00,
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-12-31T23:59:59Z",
  "location": "Rural Village, State",
  "category": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Education"
  },
  "organizer": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "Education Foundation",
    "verified": true
  },
  "urgency": 8,
  "status": "active",
  "beneficiary_name": "Rural School Foundation",
  "beneficiary_type": "organization",
  "beneficiary_count": 500,
  "payment_methods": [
    {
      "id": 1,
      "payment_method_id": 1,
      "code": "transfer",
      "name": "Bank Transfer",
      "instructions": "Transfer to account 123456"
    }
  ]
}
```

#### List Campaigns (Admin Only)
```http
GET /api/campaigns
Authorization: Bearer {token}

Response: 200 OK
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440001",
    "title": "Help Build School in Rural Area",
    "goal": 50000.00,
    "status": "active"
  }
]
```

#### Create Campaign (Admin Only)
```http
POST /api/campaigns
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "New Campaign",
  "description": "Campaign description",
  "goal": 10000.00,
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-12-31T23:59:59Z",
  "location": "City, State",
  "category": "550e8400-e29b-41d4-a716-446655440001",
  "organizer_id": "660e8400-e29b-41d4-a716-446655440001",
  "urgency": 5,
  "beneficiary_name": "Local Community",
  "beneficiary_type": "community",
  "beneficiary_count": 1000
}

Response: 201 Created
{
  "id": "new-campaign-uuid",
  "title": "New Campaign",
  ...
}
```

#### Update Campaign (Admin or Owner)
```http
PUT /api/campaigns/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Updated Campaign Title",
  "goal": 15000.00
}

Response: 200 OK
{
  "id": "770e8400-e29b-41d4-a716-446655440001",
  "title": "Updated Campaign Title",
  ...
}
```

#### Delete Campaign (Admin Only)
```http
DELETE /api/campaigns/{id}
Authorization: Bearer {token}

Response: 204 No Content
```

#### Upload Campaign Image (Admin Only)
```http
POST /api/campaigns/{id}/upload
Authorization: Bearer {token}
Content-Type: multipart/form-data

file: [binary image data]

Response: 200 OK
{
  "url": "https://s3.amazonaws.com/bucket/campaigns/uuid/image.jpg",
  "key": "campaigns/uuid/image.jpg",
  "file_name": "image.jpg",
  "size": 102400
}
```

### Campaign Receipts Endpoints

#### Get Campaign Receipts (Public)
```http
GET /api/campaigns/{campaignId}/receipts

Response: 200 OK
[
  {
    "id": "receipt-uuid",
    "campaign_id": "campaign-uuid",
    "provider": "Office Supply Store",
    "name": "School Supplies Purchase",
    "description": "Notebooks, pencils, and other supplies for students",
    "total": 1500.00,
    "quantity": 50,
    "date": "2024-06-15T10:00:00Z",
    "document_url": "https://s3.amazonaws.com/bucket/receipts/uuid/document.pdf",
    "note": "Purchased for initial distribution"
  }
]
```

#### Get Receipt by ID (Public)
```http
GET /api/campaigns/{campaignId}/receipts/{id}

Response: 200 OK
{
  "id": "receipt-uuid",
  "campaign_id": "campaign-uuid",
  "provider": "Office Supply Store",
  "name": "School Supplies Purchase",
  "description": "Notebooks, pencils, and other supplies for students",
  "total": 1500.00,
  "quantity": 50,
  "date": "2024-06-15T10:00:00Z",
  "document_url": "https://s3.amazonaws.com/bucket/receipts/uuid/document.pdf",
  "note": "Purchased for initial distribution"
}
```

#### Create Receipt (Admin Only)
```http
POST /api/campaigns/{campaignId}/receipts
Authorization: Bearer {token}
Content-Type: application/json

{
  "provider": "Medical Supply Company",
  "name": "Medical Equipment Purchase",
  "description": "Essential medical supplies for clinic",
  "total": 5000.00,
  "quantity": 10,
  "date": "2024-06-20T15:00:00Z",
  "note": "Emergency supplies for new clinic"
}

Response: 201 Created
{
  "id": "new-receipt-uuid",
  "campaign_id": "campaign-uuid",
  ...
}
```

#### Update Receipt (Admin Only)
```http
PUT /api/campaigns/{campaignId}/receipts/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "total": 5500.00,
  "note": "Updated amount after additional items"
}

Response: 200 OK
{
  "id": "receipt-uuid",
  "total": 5500.00,
  ...
}
```

#### Delete Receipt (Admin Only)
```http
DELETE /api/campaigns/{campaignId}/receipts/{id}
Authorization: Bearer {token}

Response: 204 No Content
```

#### Upload Receipt Document (Admin Only)
```http
POST /api/campaigns/{campaignId}/receipts/{id}/upload
Authorization: Bearer {token}
Content-Type: multipart/form-data

file: [binary PDF data]

Response: 200 OK
{
  "url": "https://s3.amazonaws.com/bucket/receipts/uuid/document.pdf",
  "key": "receipts/uuid/document.pdf",
  "file_name": "receipt.pdf",
  "size": 204800
}
```

### Campaign Activity Endpoints

#### Get Campaign Activities
```http
GET /api/campaigns/{campaignId}/activities

Response: 200 OK
[
  {
    "id": "activity-uuid",
    "campaign_id": "campaign-uuid",
    "title": "Campaign Update",
    "description": "We've reached 50% of our goal!",
    "date": "2024-06-15T10:00:00Z",
    "type": "update",
    "author": "John Doe"
  }
]
```

#### Create Activity (Admin Only)
```http
POST /api/campaigns/{campaignId}/activities
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Milestone Reached",
  "description": "We've completed the first phase",
  "date": "2024-06-20T15:00:00Z",
  "type": "milestone",
  "author": "Campaign Manager"
}

Response: 201 Created
{
  "id": "new-activity-uuid",
  "campaign_id": "campaign-uuid",
  ...
}
```

#### Update Activity (Admin Only)
```http
PUT /api/campaigns/{campaignId}/activities/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Updated Title",
  "description": "Updated description"
}

Response: 200 OK
{
  "id": "activity-uuid",
  "title": "Updated Title",
  ...
}
```

#### Delete Activity (Admin Only)
```http
DELETE /api/campaigns/{campaignId}/activities/{id}
Authorization: Bearer {token}

Response: 204 No Content
```

### Donation Endpoints

#### Get Campaign Donations (Public)
```http
GET /api/campaigns/{campaignId}/donations

Response: 200 OK
[
  {
    "id": "donation-uuid",
    "campaign_id": "campaign-uuid",
    "amount": 100.00,
    "donor_id": "donor-uuid",
    "date": "2024-06-01T12:00:00Z",
    "message": "Happy to help!",
    "is_anonymous": false,
    "payment_method": {
      "id": 1,
      "code": "transfer",
      "name": "Bank Transfer"
    },
    "status": "completed"
  }
]
```

#### Get Donation by ID (Public)
```http
GET /api/campaigns/{campaignId}/donations/{id}

Response: 200 OK
{
  "id": "donation-uuid",
  "campaign_id": "campaign-uuid",
  "amount": 100.00,
  "donor_id": "donor-uuid",
  "date": "2024-06-01T12:00:00Z",
  "message": "Happy to help!",
  "is_anonymous": false,
  "payment_method": {
    "id": 1,
    "code": "transfer",
    "name": "Bank Transfer"
  },
  "status": "completed"
}
```

#### Create Donation (Admin Only)
```http
POST /api/campaigns/{campaignId}/donations
Authorization: Bearer {token}
Content-Type: application/json

{
  "amount": 250.00,
  "donor_id": "donor-uuid",
  "message": "Keep up the good work!",
  "is_anonymous": false,
  "payment_method_id": 1
}

Response: 201 Created
{
  "id": "new-donation-uuid",
  "campaign_id": "campaign-uuid",
  "amount": 250.00,
  "date": "2024-06-01T12:00:00Z",
  "donor_id": "donor-uuid",
  "message": "Keep up the good work!",
  "is_anonymous": false,
  "payment_method_id": 1,
  "status": "pending"
}
```

#### Update Donation (Admin Only)
```http
PUT /api/campaigns/{campaignId}/donations/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": "completed",
  "message": "Updated message"
}

Response: 200 OK
{
  "id": "donation-uuid",
  "status": "completed",
  "message": "Updated message",
  ...
}
```

### Donor Endpoints

#### List Donors
```http
GET /api/donors
Authorization: Bearer {token}

Response: 200 OK
[
  {
    "id": "donor-uuid",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "is_verified": true
  }
]
```

#### Get Donor by ID
```http
GET /api/donors/{id}
Authorization: Bearer {token}

Response: 200 OK
{
  "id": "donor-uuid",
  "first_name": "John",
  "last_name": "Doe",
  ...
}
```

#### Create Donor
```http
POST /api/donors
Authorization: Bearer {token}
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane@example.com",
  "phone": "+0987654321"
}

Response: 201 Created
{
  "id": "new-donor-uuid",
  "first_name": "Jane",
  ...
}
```

#### Update Donor
```http
PUT /api/donors/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "phone": "+1111111111",
  "is_verified": true
}

Response: 200 OK
{
  "id": "donor-uuid",
  "phone": "+1111111111",
  ...
}
```

### Organizer Endpoints

#### List Organizers
```http
GET /api/organizers
Authorization: Bearer {token}

Response: 200 OK
[
  {
    "id": "organizer-uuid",
    "user_id": "user-uuid",
    "name": "Education Foundation",
    "avatar": "https://example.com/avatar.jpg",
    "verified": true,
    "email": "contact@foundation.org",
    "phone": "+1234567890",
    "website": "https://foundation.org",
    "description": "Non-profit education organization"
  }
]
```

#### Get Organizer by ID
```http
GET /api/organizers/{id}
Authorization: Bearer {token}

Response: 200 OK
{
  "id": "organizer-uuid",
  "name": "Education Foundation",
  ...
}
```

#### Create Organizer
```http
POST /api/organizers
Authorization: Bearer {token}
Content-Type: application/json

{
  "user_id": "user-uuid",
  "name": "New Organization",
  "email": "contact@neworg.com",
  "phone": "+1234567890",
  "website": "https://neworg.com",
  "description": "New organization description"
}

Response: 201 Created
{
  "id": "new-organizer-uuid",
  "name": "New Organization",
  ...
}
```

#### Update Organizer
```http
PUT /api/organizers/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "verified": true,
  "description": "Updated description"
}

Response: 200 OK
{
  "id": "organizer-uuid",
  "verified": true,
  ...
}
```

### Campaign Category Endpoints

#### List Categories
```http
GET /api/categories

Response: 200 OK
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Education",
    "description": "Campaigns related to education and learning"
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "name": "Health",
    "description": "Medical and health-related campaigns"
  }
]
```

#### Get Category by ID
```http
GET /api/categories/{id}

Response: 200 OK
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "Education",
  "description": "Campaigns related to education and learning"
}
```

#### Create Category (Admin Only)
```http
POST /api/categories
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Technology",
  "description": "Technology and innovation campaigns"
}

Response: 201 Created
{
  "id": "new-category-uuid",
  "name": "Technology",
  "description": "Technology and innovation campaigns"
}
```

#### Update Category (Admin Only)
```http
PUT /api/categories/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "description": "Updated description for category"
}

Response: 200 OK
{
  "id": "category-uuid",
  "name": "Education",
  "description": "Updated description for category"
}
```

### Payment Method Endpoints

#### List Payment Methods
```http
GET /api/payment-methods

Response: 200 OK
[
  {
    "id": 1,
    "code": "transfer",
    "name": "Bank Transfer",
    "is_active": true
  },
  {
    "id": 2,
    "code": "cash",
    "name": "Cash",
    "is_active": true
  }
]
```

#### Create Campaign Payment Method
```http
POST /api/campaign-payment-methods
Authorization: Bearer {token}
Content-Type: application/json

{
  "campaign_id": "campaign-uuid",
  "payment_method_id": 1,
  "instructions": "Transfer to account 123456",
  "transfer_details": {
    "bank_name": "Example Bank",
    "account_holder": "Foundation Name",
    "account_number": "123456789",
    "routing_number": "987654321",
    "swift_code": "EXMBUS33",
    "additional_info": "Reference: Campaign Name"
  }
}

Response: 201 Created
{
  "id": 1,
  "campaign_id": "campaign-uuid",
  "payment_method_id": 1,
  ...
}
```

#### Create Cash Location
```http
POST /api/campaign-payment-methods/{id}/cash-locations
Authorization: Bearer {token}
Content-Type: application/json

{
  "location_name": "Main Office",
  "address": "123 Main St, City, State",
  "contact_person": "John Doe",
  "contact_phone": "+1234567890",
  "available_hours": "Mon-Fri 9AM-5PM"
}

Response: 201 Created
{
  "id": 1,
  "campaign_payment_method_id": 1,
  "location_name": "Main Office",
  ...
}
```

### RBAC Endpoints

#### List Roles
```http
GET /api/roles
Authorization: Bearer {token}

Response: 200 OK
[
  {
    "id": "11111111-1111-1111-1111-111111111111",
    "name": "admin",
    "description": "Administrator with full system access",
    "is_active": true
  },
  {
    "id": "22222222-2222-2222-2222-222222222222",
    "name": "donor",
    "description": "Donor who can create and manage donations",
    "is_active": true
  }
]
```

#### List Permissions
```http
GET /api/permissions
Authorization: Bearer {token}

Response: 200 OK
[
  {
    "id": "permission-uuid",
    "name": "campaigns:create",
    "resource": "campaigns",
    "action": "create",
    "description": "Create new campaigns"
  }
]
```

---

## Database Schemas

### Core Tables

#### users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role_id UUID NOT NULL REFERENCES roles(id),
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    reset_token VARCHAR(255),
    reset_token_expires_at TIMESTAMP WITH TIME ZONE,
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### campaigns
```sql
CREATE TABLE campaigns (
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
    beneficiary_name VARCHAR(255),
    beneficiary_type VARCHAR(50),
    beneficiary_count INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_end_date_after_start CHECK (end_date > start_date)
);
```

#### donations
```sql
CREATE TABLE donations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id),
    donor_id UUID NOT NULL REFERENCES donors(id),
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    message TEXT,
    is_anonymous BOOLEAN DEFAULT false,
    payment_method_id INTEGER REFERENCES payment_methods(id),
    status donation_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Enum type for donation status
CREATE TYPE donation_status AS ENUM ('completed', 'pending', 'failed', 'refunded');
```

#### organizers
```sql
CREATE TABLE organizers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    avatar VARCHAR(500),
    verified BOOLEAN DEFAULT FALSE,
    email VARCHAR(255),
    phone VARCHAR(20),
    website VARCHAR(500),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### donors
```sql
CREATE TABLE donors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT false,
    phone VARCHAR(50),
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### campaign_categories
```sql
CREATE TABLE campaign_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### activities
```sql
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    type VARCHAR(50),
    author VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### receipts
```sql
CREATE TABLE receipts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    provider VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total DECIMAL(10,2) NOT NULL CHECK (total > 0),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    document_url VARCHAR(500),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### RBAC Tables

#### roles
```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Default roles with fixed UUIDs
INSERT INTO roles (id, name, description) VALUES 
    ('11111111-1111-1111-1111-111111111111', 'admin', 'Administrator with full system access'),
    ('22222222-2222-2222-2222-222222222222', 'donor', 'Donor who can create and manage donations'),
    ('33333333-3333-3333-3333-333333333333', 'guest', 'Guest with read-only access to public content');
```

#### permissions
```sql
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(20) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Permissions format: resource:action
-- Resources: campaigns, donations, users, categories, organizers, donors
-- Actions: create, read, update, delete
```

#### role_permissions
```sql
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);
```

### Payment System Tables

#### payment_methods
```sql
CREATE TABLE payment_methods (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Default payment methods
INSERT INTO payment_methods (code, name) VALUES
    ('transfer', 'Bank Transfer'),
    ('cash', 'Cash'),
    ('mercadopago', 'MercadoPago'),
    ('paypal', 'PayPal'),
    ('crypto', 'Cryptocurrency');
```

#### campaign_payment_methods
```sql
CREATE TABLE campaign_payment_methods (
    id SERIAL PRIMARY KEY,
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    payment_method_id INTEGER NOT NULL REFERENCES payment_methods(id),
    instructions TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(campaign_id, payment_method_id)
);
```

#### transfer_details
```sql
CREATE TABLE transfer_details (
    id SERIAL PRIMARY KEY,
    campaign_payment_method_id INTEGER NOT NULL REFERENCES campaign_payment_methods(id) ON DELETE CASCADE,
    bank_name VARCHAR(255) NOT NULL,
    account_holder VARCHAR(255) NOT NULL,
    account_number VARCHAR(100) NOT NULL,
    routing_number VARCHAR(50),
    swift_code VARCHAR(20),
    additional_info TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### cash_locations
```sql
CREATE TABLE cash_locations (
    id SERIAL PRIMARY KEY,
    campaign_payment_method_id INTEGER NOT NULL REFERENCES campaign_payment_methods(id) ON DELETE CASCADE,
    location_name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    contact_person VARCHAR(255),
    contact_phone VARCHAR(50),
    available_hours VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Database Relationships

```mermaid
erDiagram
    users ||--o{ organizers : "has"
    users ||--|| roles : "belongs to"
    roles ||--o{ role_permissions : "has"
    permissions ||--o{ role_permissions : "has"
    
    organizers ||--o{ campaigns : "creates"
    campaign_categories ||--o{ campaigns : "categorizes"
    campaigns ||--o{ donations : "receives"
    campaigns ||--o{ activities : "has"
    campaigns ||--o{ receipts : "has"
    campaigns ||--o{ campaign_payment_methods : "has"
    
    donors ||--o{ donations : "makes"
    payment_methods ||--o{ donations : "used in"
    payment_methods ||--o{ campaign_payment_methods : "configures"
    
    campaign_payment_methods ||--o{ transfer_details : "has"
    campaign_payment_methods ||--o{ cash_locations : "has"
```

---

## Authentication Flow

### 1. User Registration Flow
```
1. User submits registration form with email, password, first_name, last_name
2. System validates input data
3. Password is hashed using bcrypt
4. User record created with default role (guest)
5. Organizer record created linked to user
6. User ID returned to client
```

### 2. Login Flow
```
1. User submits email and password
2. System validates credentials against database
3. If valid, JWT token is generated containing:
   - sub: user_id
   - role_id: user's role UUID
   - role: role name (admin/donor/guest)
   - exp: expiration timestamp (24 hours)
4. Token returned to client with token_type "Bearer"
5. Client includes token in Authorization header for subsequent requests
```

### 3. JWT Token Structure
```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000",  // user_id
  "role_id": "11111111-1111-1111-1111-111111111111",  // role UUID
  "role": "admin",  // role name
  "exp": 1719590400,  // expiration timestamp
  "iat": 1719504000   // issued at timestamp
}
```

### 4. Request Authentication
```
1. Client sends request with Authorization header:
   Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   
2. RequireAuth middleware:
   - Extracts token from header
   - Validates token signature (HS256)
   - Checks token expiration
   - Extracts claims and stores in context:
     - user_id
     - role_id
     - role (name)
     
3. Request proceeds to handler with user context
```

### 5. RBAC Authorization Flow
```
1. After authentication, RBAC middleware checks:
   a. Role-based: Does user's role match required role?
   b. Permission-based: Does user's role have required permission?
   c. Ownership-based: Is user the owner of the resource?

2. Permission check process:
   - Get user's role from context
   - Query role_permissions for role's permissions
   - Check if required permission exists
   - Allow/deny access accordingly

3. Ownership check (for user-specific resources):
   - Extract resource ID from request
   - Query database for resource owner
   - Compare with authenticated user_id
   - Allow if match or user is admin
```

### 6. Password Reset Flow
```
1. User requests password reset with email
2. System generates unique reset token
3. Token stored in users.reset_token with expiration
4. Email sent with reset link (not implemented)
5. User submits new password with token
6. System validates token and expiration
7. Password updated and token cleared
```

### 7. Permission Structure
Permissions follow the format: `resource:action`

**Resources:**
- campaigns
- donations
- users
- categories
- organizers
- donors

**Actions:**
- create
- read
- update
- delete

**Example Permission Matrix:**

| Role  | Permissions |
|-------|------------|
| admin | All permissions (full CRUD on all resources) |
| donor | campaigns:read, donations:*, donors:*, users:update (own) |
| guest | campaigns:read, categories:read, organizers:read |

### 8. Security Considerations
- Passwords hashed with bcrypt (cost factor 10)
- JWT tokens expire after 24 hours
- Reset tokens have expiration time
- HTTPS required in production
- CORS configured for frontend origins
- SQL injection prevented via parameterized queries (GORM)
- Input validation on all endpoints
- Rate limiting recommended (not implemented)

---

## Business Logic

### Campaign Management
*[Placeholder for your business logic documentation]*

### Donation Processing
*[Placeholder for your business logic documentation]*

### Payment Workflows
*[Placeholder for your business logic documentation]*

### Notification System
*[Placeholder for your business logic documentation]*

### Reporting & Analytics
*[Placeholder for your business logic documentation]*

### Data Validation Rules
*[Placeholder for your business logic documentation]*

### Integration Points
*[Placeholder for your business logic documentation]*

---

## Environment Variables

Required environment variables for the application:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=microservice_user
DB_PASSWORD=microservice_password
DB_NAME=microservice_db
DB_SSLMODE=disable

# API Configuration
API_PORT=9999

# Docker Configuration
DB_PORT_EXTERNAL=5432  # 5440 for development

# JWT Configuration
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h

# AWS S3 Configuration (for file uploads)
AWS_REGION=us-east-1
AWS_S3_BUCKET=dona-tutti-files
AWS_ACCESS_KEY_ID=your-aws-access-key-id
AWS_SECRET_ACCESS_KEY=your-aws-secret-access-key
```

---

## Quick Start Commands

```bash
# Development
make build-dev   # Build and start with hot reload
make start-dev   # Start without rebuild

# Production
make build-prod  # Build and start production
make start-prod  # Start without rebuild

# Database Access
docker-compose exec postgres psql -U microservice_user -d microservice_db

# View Logs
make logs        # All services
make logs-api    # API only
make logs-db     # Database only

# Stop Services
make stop        # Stop all containers
make clean       # Stop and remove volumes
```

---

## API Testing

Use the provided `test_api.sh` script or the curl examples above to test endpoints.

For authenticated endpoints, first obtain a token via login:
```bash
TOKEN=$(curl -s -X POST http://localhost:9999/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | jq -r '.token.access_token')

# Use token in subsequent requests
curl -H "Authorization: Bearer $TOKEN" http://localhost:9999/api/users/me
```

---

*Document generated for Dona Tutti API v1.0*
*Last updated: 2024*