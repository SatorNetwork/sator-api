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
	role, err := jwt.RoleFromContext(ctx)
	if err != nil {
		role = RoleGuest.String()
	} else {
		if role == "" {
			role = RoleUser.String()
		}
	}

	return CheckRole(Role(role), allowedRoles...)
}
