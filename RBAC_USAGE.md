# RBAC System Usage Guide

## Overview

This guide explains how to use the Role-Based Access Control (RBAC) system implemented in the Dona Tutti API.

## Architecture

The RBAC system consists of:

- **Roles**: admin, donor, guest
- **Permissions**: resource:action combinations (e.g., "campaigns:create")
- **Resources**: campaigns, donations, users, categories, organizers, donors
- **Actions**: create, read, update, delete

## Default Role Permissions

### Admin Role
- Full access to all resources and actions
- Can manage users, roles, and permissions
- Can approve/reject campaigns
- Access to financial reports

### Donor Role
- Create, read, update, delete own donations
- Read access to campaigns and categories
- Update own user profile
- Create donor profiles

### Guest Role
- Read-only access to public campaigns
- Read access to categories and organizers
- Cannot create or modify any resources

## Middleware Usage

### 1. Basic Role-Based Authorization

```go
// Require specific role
rbacMiddleware.RequireRole("admin")
rbacMiddleware.RequireRole("donor")

// Require any of multiple roles
rbacMiddleware.RequireAnyRole("admin", "donor")

// Require all specified roles (rarely used)
rbacMiddleware.RequireAllRoles("admin", "donor")
```

### 2. Permission-Based Authorization

```go
// Require specific permission
rbacMiddleware.RequirePermission("campaigns:create")
rbacMiddleware.RequirePermission("donations:update")

// Require any of multiple permissions
rbacMiddleware.RequirePermissionWithConfig(middleware.PermissionConfig{
    Permissions: []string{"campaigns:create", "campaigns:update"},
    RequireAll:  false,
})

// Require all specified permissions
rbacMiddleware.RequireAllPermissions("campaigns:create", "campaigns:update")
```

### 3. Ownership-Based Authorization

```go
// Basic ownership check (user can access their own resources)
rbacMiddleware.RequireOwnership()

// Advanced ownership configuration
rbacMiddleware.RequireOwnershipWithConfig(middleware.OwnershipConfig{
    Resource:         "donations",
    ResourceIDParam:  "id",
    AllowAdminBypass: true,
})
```

### 4. Combined Authorization (OR Logic)

```go
// User must be admin OR own the resource
rbacMiddleware.Combine(
    rbacMiddleware.RequireRole("admin"),
    rbacMiddleware.RequireOwnership(),
)

// User must have permission OR be admin
rbacMiddleware.Combine(
    rbacMiddleware.RequirePermission("campaigns:update"),
    rbacMiddleware.RequireRole("admin"),
)
```

## Implementation Examples

### Campaign Routes with RBAC

```go
func RegisterRoutes(g *echo.Group, service Service, db *gorm.DB) {
    handler := NewHandler(service)
    rbacMiddleware := middleware.NewRBACMiddleware(db)
    
    campaignGroup := g.Group("/campaigns")
    
    // Public routes
    campaignGroup.GET("", handler.ListCampaigns)
    campaignGroup.GET("/:id", handler.GetCampaign)
    
    // Protected routes
    authGroup := campaignGroup.Group("", middleware.RequireAuth())
    
    // Admin-only routes
    authGroup.POST("", 
        rbacMiddleware.RequireRole("admin"), 
        handler.CreateCampaign)
    authGroup.DELETE("/:id", 
        rbacMiddleware.RequireRole("admin"), 
        handler.DeleteCampaign)
    
    // Admin or owner routes
    authGroup.PUT("/:id", 
        rbacMiddleware.Combine(
            rbacMiddleware.RequireRole("admin"),
            rbacMiddleware.RequireOwnership(),
        ), 
        handler.UpdateCampaign)
}
```

### User Routes with RBAC

```go
func RegisterRoutes(g *echo.Group, service Service, db *gorm.DB) {
    handler := NewHandler(service)
    rbacMiddleware := middleware.NewRBACMiddleware(db)
    
    userGroup := g.Group("/users", middleware.RequireAuth())
    
    // Admin can list all users
    userGroup.GET("", 
        rbacMiddleware.RequireRole("admin"), 
        handler.ListUsers)
    
    // Users can view their own profile, admins can view any
    userGroup.GET("/:id", 
        rbacMiddleware.Combine(
            rbacMiddleware.RequireRole("admin"),
            rbacMiddleware.RequireOwnership(),
        ), 
        handler.GetUser)
    
    // Users can update their own profile only
    userGroup.PUT("/:id", 
        rbacMiddleware.RequireOwnership(), 
        handler.UpdateUser)
}
```

## API Endpoints for RBAC Management

### Check User Permissions

```bash
GET /api/auth/check-permission?permission=campaigns:create
Authorization: Bearer <token>
```

### Get User Context

```bash
GET /api/auth/user-context
Authorization: Bearer <token>
```

### List Roles (Admin only)

```bash
GET /api/roles
Authorization: Bearer <admin-token>
```

### List Permissions (Admin only)

```bash
GET /api/permissions
Authorization: Bearer <admin-token>
```

## Database Setup

1. Run the RBAC migration:
```bash
# The migration is automatically applied on server startup
# File: migrations/20240325000000_rbac_system.sql
```

2. Default roles and permissions are seeded automatically.

## JWT Token Structure

The JWT token now includes role information:

```json
{
  "sub": "user-uuid",
  "role_id": "role-uuid", 
  "role": "admin",
  "exp": 1234567890
}
```

## Migration Guide

To migrate existing routes to use RBAC:

1. Add database parameter to RegisterRoutes function
2. Initialize RBAC middleware in the function
3. Wrap routes with appropriate authorization middleware
4. Test with different user roles

Example:
```go
// Before
func RegisterRoutes(g *echo.Group, service Service) {
    // routes without authorization
}

// After  
func RegisterRoutes(g *echo.Group, service Service, db *gorm.DB) {
    rbacMiddleware := middleware.NewRBACMiddleware(db)
    // routes with authorization
}
```

## Testing RBAC

1. Create users with different roles using the registration API
2. Login to get JWT tokens for each role
3. Test endpoints with different tokens to verify authorization
4. Use the permission check endpoint to verify permissions

## Common Patterns

### Read-Only for Guests, Full Access for Admin

```go
// Public read
group.GET("/:id", handler.Get)

// Admin write operations
authGroup := group.Group("", middleware.RequireAuth())
authGroup.POST("", rbacMiddleware.RequireRole("admin"), handler.Create)
authGroup.PUT("/:id", rbacMiddleware.RequireRole("admin"), handler.Update)
authGroup.DELETE("/:id", rbacMiddleware.RequireRole("admin"), handler.Delete)
```

### Owner or Admin Access

```go
authGroup.PUT("/:id", 
    rbacMiddleware.Combine(
        rbacMiddleware.RequireRole("admin"),
        rbacMiddleware.RequireOwnership(),
    ), 
    handler.Update)
```

### Permission-Based with Ownership Fallback

```go
authGroup.PUT("/:id", 
    rbacMiddleware.RequirePermissionWithConfig(middleware.PermissionConfig{
        Permissions:    []string{"resource:update"},
        AllowOwnership: true,
    }), 
    handler.Update)
```

This RBAC system provides flexible, secure authorization that can be easily extended as your application grows.