package rbac

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Service interface {
	// Permission checking
	HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
	HasRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error)
	HasAnyRole(ctx context.Context, userID uuid.UUID, roleNames []string) (bool, error)

	// Context operations
	GetUserAuthContext(ctx context.Context, userID uuid.UUID) (interface{}, error)

	// Resource ownership validation
	ValidateResourceOwnership(ctx context.Context, userID uuid.UUID, resource string, resourceID uuid.UUID) (bool, error)

	// Role management
	ListRoles(ctx context.Context) ([]Role, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)

	// Permission management
	ListPermissions(ctx context.Context) ([]Permission, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	authCtxInterface, err := s.repo.GetUserAuthContext(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user auth context: %w", err)
	}

	// Type assert to LocalAuthContext
	authCtx, ok := authCtxInterface.(*LocalAuthContext)
	if !ok {
		return false, fmt.Errorf("invalid auth context type")
	}

	// Check if user has the specific permission
	for _, userPerm := range authCtx.Permissions {
		if userPerm == permission {
			return true, nil
		}
	}

	return false, nil
}

func (s *service) HasRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error) {
	authCtxInterface, err := s.repo.GetUserAuthContext(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user auth context: %w", err)
	}

	// Type assert to LocalAuthContext
	authCtx, ok := authCtxInterface.(*LocalAuthContext)
	if !ok {
		return false, fmt.Errorf("invalid auth context type")
	}

	return strings.EqualFold(authCtx.RoleName, roleName), nil
}

func (s *service) HasAnyRole(ctx context.Context, userID uuid.UUID, roleNames []string) (bool, error) {
	authCtxInterface, err := s.repo.GetUserAuthContext(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user auth context: %w", err)
	}

	// Type assert to LocalAuthContext
	authCtx, ok := authCtxInterface.(*LocalAuthContext)
	if !ok {
		return false, fmt.Errorf("invalid auth context type")
	}

	for _, roleName := range roleNames {
		if strings.EqualFold(authCtx.RoleName, roleName) {
			return true, nil
		}
	}

	return false, nil
}

func (s *service) GetUserAuthContext(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	authCtx, err := s.repo.GetUserAuthContext(ctx, userID)
	return authCtx, err
}

func (s *service) ValidateResourceOwnership(ctx context.Context, userID uuid.UUID, resource string, resourceID uuid.UUID) (bool, error) {
	// This method should validate if the user owns the specific resource
	// Implementation depends on your business logic
	// For now, we'll implement basic ownership validation

	switch resource {
	case ResourceUsers:
		// Users can only access their own profile
		return userID == resourceID, nil
	case ResourceDonations:
		// TODO: Check if user owns the donation
		// This would require a query to the donations table
		return s.validateDonationOwnership(ctx, userID, resourceID)
	case ResourceDonors:
		// TODO: Check if user owns the donor profile
		// This would require a query to the donors table
		return s.validateDonorOwnership(ctx, userID, resourceID)
	default:
		// For other resources, we don't have ownership validation yet
		return false, nil
	}
}

func (s *service) validateDonationOwnership(ctx context.Context, userID uuid.UUID, donationID uuid.UUID) (bool, error) {
	// TODO: Implement donation ownership validation
	// This would require access to the donation repository
	// For now, return false
	return false, nil
}

func (s *service) validateDonorOwnership(ctx context.Context, userID uuid.UUID, donorID uuid.UUID) (bool, error) {
	// TODO: Implement donor ownership validation
	// This would require access to the donor repository
	// For now, return false
	return false, nil
}

func (s *service) ListRoles(ctx context.Context) ([]Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *service) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	return s.repo.GetRoleByName(ctx, name)
}

func (s *service) ListPermissions(ctx context.Context) ([]Permission, error) {
	return s.repo.ListPermissions(ctx)
}
