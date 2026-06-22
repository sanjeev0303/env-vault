package domain

import "time"

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleEditor Role = "editor"
	RoleViewer Role = "viewer"
)

type Permission string

const (
	PermCreateSecret  Permission = "CREATE_SECRET"
	PermUpdateSecret  Permission = "UPDATE_SECRET"
	PermDeleteSecret  Permission = "DELETE_SECRET"
	PermViewSecret    Permission = "VIEW_SECRET"
	PermManageMembers Permission = "MANAGE_MEMBERS"
	PermExportSecrets Permission = "EXPORT_SECRETS"
	PermManageProject Permission = "MANAGE_PROJECT"
)

// RolePermissions defines which permissions each role has.
var RolePermissions = map[Role][]Permission{
	RoleOwner: {
		PermCreateSecret, PermUpdateSecret, PermDeleteSecret, PermViewSecret,
		PermManageMembers, PermExportSecrets, PermManageProject,
	},
	RoleAdmin: {
		PermCreateSecret, PermUpdateSecret, PermDeleteSecret, PermViewSecret,
		PermManageMembers, PermExportSecrets,
	},
	RoleEditor: {
		PermCreateSecret, PermUpdateSecret, PermDeleteSecret, PermViewSecret,
		PermExportSecrets,
	},
	RoleViewer: {
		PermViewSecret,
	},
}

// HasPermission checks if a role has a specific permission.
func HasPermission(role Role, perm Permission) bool {
	perms, ok := RolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}

// ValidRole checks if a string is a valid role.
func ValidRole(r string) bool {
	switch Role(r) {
	case RoleOwner, RoleAdmin, RoleEditor, RoleViewer:
		return true
	}
	return false
}

type ProjectMember struct {
	ID        string
	ProjectID string
	UserID    string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
