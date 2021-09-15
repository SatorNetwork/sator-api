package rbac

// Predefined user roles.
const (
	AvailableForAllRoles        Role = "available_for_all_roles"
	AvailableForAuthorizedUsers Role = "available_for_authorized_users"

	RoleAdmin          Role = "admin"
	RoleContentManager Role = "content_manager"
	RoleGuest          Role = "guest"
	RoleShowRunner     Role = "show_runner"
	RoleUser           Role = "user"
)

type Role string

func (r Role) String() string {
	return string(r)
}
