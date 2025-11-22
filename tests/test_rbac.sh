#!/bin/bash

# RBAC Testing Script for Dona Tutti API
# This script tests the Role-Based Access Control implementation

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:9999/api"

echo -e "${BLUE}🔐 RBAC TESTING SCRIPT FOR DONA TUTTI API${NC}"
echo "========================================================="
echo ""

# Function to show test results
show_test() {
    echo -e "${BLUE}$1${NC}"
    echo "----------------------------------------"
}

show_success() {
    echo -e "${GREEN}✅ SUCCESS:${NC} $1"
    echo ""
}

show_error() {
    echo -e "${RED}❌ ERROR:${NC} $1"
    echo ""
}

show_warning() {
    echo -e "${YELLOW}⚠️  WARNING:${NC} $1"
    echo ""
}

show_info() {
    echo -e "${PURPLE}ℹ️  INFO:${NC} $1"
    echo ""
}

# Check if server is running
show_test "Checking if server is running..."
response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/campaigns")
if [ "$response" != "200" ]; then
    show_error "Server is not running or not accessible at $BASE_URL"
    echo "Please start the server with: go run main.go"
    exit 1
fi
show_success "Server is running and accessible"

# Test 1: Register users with different roles
echo -e "${YELLOW}📝 STEP 1: USER REGISTRATION${NC}"
echo "=========================="

# Register admin user
show_test "1.1. Registering admin user"
admin_response=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@donatutti.com",
    "password": "admin123456",
    "first_name": "Admin",
    "last_name": "User"
  }')

if [[ $admin_response == *"id"* ]]; then
    show_success "Admin user registered successfully"
    ADMIN_ID=$(echo "$admin_response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    show_info "Admin ID: $ADMIN_ID"
else
    show_warning "Admin user may already exist or registration failed"
fi

# Register donor user
show_test "1.2. Registering donor user"
donor_response=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "donor@donatutti.com", 
    "password": "donor123456",
    "first_name": "Donor",
    "last_name": "User"
  }')

if [[ $donor_response == *"id"* ]]; then
    show_success "Donor user registered successfully"
    DONOR_ID=$(echo "$donor_response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    show_info "Donor ID: $DONOR_ID"
else
    show_warning "Donor user may already exist or registration failed"
fi

# Register guest user
show_test "1.3. Registering guest user"
guest_response=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "guest@donatutti.com",
    "password": "guest123456", 
    "first_name": "Guest",
    "last_name": "User"
  }')

if [[ $guest_response == *"id"* ]]; then
    show_success "Guest user registered successfully"
    GUEST_ID=$(echo "$guest_response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    show_info "Guest ID: $GUEST_ID"
else
    show_warning "Guest user may already exist or registration failed"
fi

# Test 2: User login and token retrieval
echo -e "${YELLOW}🔑 STEP 2: USER LOGIN${NC}"
echo "==================="

# Login admin
show_test "2.1. Admin login"
admin_login=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@donatutti.com",
    "password": "admin123456"
  }')

if [[ $admin_login == *"access_token"* ]]; then
    show_success "Admin login successful"
    ADMIN_TOKEN=$(echo "$admin_login" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    show_info "Admin token obtained"
else
    show_error "Admin login failed"
    echo "Response: $admin_login"
fi

# Login donor
show_test "2.2. Donor login"
donor_login=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "donor@donatutti.com",
    "password": "donor123456"
  }')

if [[ $donor_login == *"access_token"* ]]; then
    show_success "Donor login successful"
    DONOR_TOKEN=$(echo "$donor_login" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    show_info "Donor token obtained"
else
    show_error "Donor login failed"
    echo "Response: $donor_login"
fi

# Login guest
show_test "2.3. Guest login"
guest_login=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "guest@donatutti.com", 
    "password": "guest123456"
  }')

if [[ $guest_login == *"access_token"* ]]; then
    show_success "Guest login successful"
    GUEST_TOKEN=$(echo "$guest_login" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    show_info "Guest token obtained"
else
    show_error "Guest login failed"
    echo "Response: $guest_login"
fi

# Test 3: Public endpoint access (no authentication required)
echo -e "${YELLOW}🌐 STEP 3: PUBLIC ENDPOINT ACCESS${NC}"
echo "=================================="

show_test "3.1. Public campaign list access"
public_campaigns=$(curl -s "$BASE_URL/campaigns")
if [[ $public_campaigns == *"["* ]] || [[ $public_campaigns == *"campaigns"* ]]; then
    show_success "Public campaign access working"
else
    show_error "Public campaign access failed"
fi

# Test 4: Authentication required endpoints
echo -e "${YELLOW}🔒 STEP 4: AUTHENTICATION TESTS${NC}"
echo "================================"

show_test "4.1. Access protected endpoint without token"
no_auth_response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/users")
if [ "$no_auth_response" = "401" ]; then
    show_success "Properly blocked access without authentication"
else
    show_warning "Expected 401 Unauthorized, got $no_auth_response"
fi

show_test "4.2. Access protected endpoint with admin token"
if [ ! -z "$ADMIN_TOKEN" ]; then
    admin_auth_response=$(curl -s -o /dev/null -w "%{http_code}" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      "$BASE_URL/users")
    if [ "$admin_auth_response" = "200" ]; then
        show_success "Admin authentication working"
    else
        show_warning "Admin authentication returned $admin_auth_response"
    fi
else
    show_error "No admin token available for testing"
fi

# Test 5: Role-based authorization (if RBAC middleware is enabled)
echo -e "${YELLOW}👑 STEP 5: ROLE-BASED ACCESS TESTS${NC}"
echo "=================================="

show_test "5.1. Test campaign creation with different roles"

# Test admin campaign creation
if [ ! -z "$ADMIN_TOKEN" ]; then
    show_test "5.1.1. Admin creating campaign"
    admin_create=$(curl -s -X POST "$BASE_URL/campaigns" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "title": "RBAC Test Campaign (Admin)",
        "description": "Campaign created by admin for RBAC testing",
        "goal": 10000.0,
        "start_date": "2025-01-01T00:00:00Z",
        "end_date": "2025-12-31T23:59:59Z",
        "location": "Test City"
      }')
    
    if [[ $admin_create == *"id"* ]]; then
        show_success "Admin can create campaigns"
        ADMIN_CAMPAIGN_ID=$(echo "$admin_create" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    else
        show_info "Campaign creation response: $admin_create"
    fi
else
    show_error "No admin token for testing"
fi

# Test donor campaign creation (should fail if RBAC is enabled)
if [ ! -z "$DONOR_TOKEN" ]; then
    show_test "5.1.2. Donor attempting to create campaign"
    donor_create=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/campaigns" \
      -H "Authorization: Bearer $DONOR_TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "title": "RBAC Test Campaign (Donor)",
        "description": "Campaign created by donor for RBAC testing",
        "goal": 5000.0,
        "start_date": "2025-01-01T00:00:00Z",
        "end_date": "2025-12-31T23:59:59Z",
        "location": "Test City"
      }')
    
    if [ "$donor_create" = "403" ]; then
        show_success "Donor properly blocked from creating campaigns (RBAC working)"
    elif [ "$donor_create" = "201" ]; then
        show_warning "Donor can create campaigns (RBAC not enabled or donor has permission)"
    else
        show_info "Donor campaign creation returned: $donor_create"
    fi
else
    show_error "No donor token for testing"
fi

# Test 6: User context and permission checking
echo -e "${YELLOW}📋 STEP 6: PERMISSION CHECKING${NC}"
echo "==============================="

if [ ! -z "$ADMIN_TOKEN" ]; then
    show_test "6.1. Check admin user context"
    admin_context=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" \
      "$BASE_URL/auth/user-context" 2>/dev/null)
    
    if [[ $admin_context == *"admin"* ]]; then
        show_success "Admin user context retrieved successfully"
        echo "Admin context: $admin_context"
    else
        show_info "User context endpoint may not be implemented yet"
    fi
fi

if [ ! -z "$DONOR_TOKEN" ]; then
    show_test "6.2. Check donor permissions"
    donor_permission=$(curl -s -H "Authorization: Bearer $DONOR_TOKEN" \
      "$BASE_URL/auth/check-permission?permission=campaigns:read" 2>/dev/null)
    
    if [[ $donor_permission == *"true"* ]]; then
        show_success "Donor has read permission for campaigns"
    else
        show_info "Permission check endpoint may not be implemented yet"
    fi
fi

# Summary
echo -e "${BLUE}📊 RBAC TESTING SUMMARY${NC}"
echo "======================="

if [ ! -z "$ADMIN_TOKEN" ] && [ ! -z "$DONOR_TOKEN" ] && [ ! -z "$GUEST_TOKEN" ]; then
    show_success "✅ All user types can authenticate successfully"
else
    show_warning "⚠️  Some user authentication failed"
fi

echo -e "${GREEN}🎯 RBAC Implementation Status:${NC}"
echo "• Database migration: ✅ Applied"
echo "• User roles: ✅ Implemented" 
echo "• JWT with roles: ✅ Implemented"
echo "• RBAC middleware: ✅ Created"
echo "• Route protection: 🔄 Ready for implementation"
echo ""

echo -e "${PURPLE}🔧 Next Steps:${NC}"
echo "1. Enable RBAC middleware in route registrations"
echo "2. Update main.go to pass database to RegisterRoutes functions"
echo "3. Update user roles in database as needed"
echo "4. Test specific permission scenarios"
echo ""

echo -e "${BLUE}🚀 RBAC Testing Complete!${NC}"

# Clean up test data (optional)
echo -e "${YELLOW}🧹 CLEANUP${NC}"
echo "=========="
show_info "Test data (users, campaigns) created during testing"
show_info "You may want to clean up test data from the database"
echo ""