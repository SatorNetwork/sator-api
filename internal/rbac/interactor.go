package rbac

import (
	"context"

	"github.com/SatorNetwork/sator-api/internal/jwt"
)

func CheckRole(userRole Role, allowedRoles ...Role) error {
	for _, v := range allowedRoles {
		switch v {
		case userRole, AvailableForAllRoles:
			return nil
		case AvailableForAuthorizedUsers:
			if userRole != RoleGuest {
				return nil
			}
		}
	}

	return ErrAccessDenied
}

func CheckRoleFromContext(ctx context.Context, allowedRoles ...Role) error {
	return CheckRole(GetRoleFromContext(ctx), allowedRoles...)
}

func GetRoleFromContext(ctx context.Context) Role {
	role, err := jwt.RoleFromContext(ctx)
	if err != nil {
		role = RoleGuest.String()
	} else {
		if role == "" {
			role = RoleUser.String()
		}
	}

	switch r := Role(role); r {
	case RoleAdmin, RoleContentManager, RoleModerator, RoleGuest, RoleShowRunner, RoleUser:
		return r
	}

	return RoleGuest
}

func IsCurrentUserHasRole(ctx context.Context, roles ...Role) bool {
	userRole := GetRoleFromContext(ctx)
	for _, r := range roles {
		if r == userRole {
			return true
		}
	}

	return false
}
