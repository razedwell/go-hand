package model

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type RoleDefinition struct {
	Role        Role
	Permissions []Permission
}

type Permission struct {
	Code *string // e.g., user.ban, user.promote, etc.
}
