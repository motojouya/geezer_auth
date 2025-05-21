package model

import (
	"time"
)

type UnsavedRole struct {
	Name        Name
	Label       Label
	Description string
}

type Role struct {
	RoleId         uint
	RegisteredDate time.Time
	UnsavedRole
}

func CreateRole(name Name, label Label, description string) UnsavedRole {
	return UnsavedRole{
		Name:        name,
		Label:       label,
		Description: description,
	}
}

func NewRole(roleId uint, name Name, label Label, description string, registeredDate time.Time) Role {
	return Role{
		RoleId:         roleId,
		RegisteredDate: registeredDate,
		UnsavedRole:    CreateRole(name, label, description),
	}
}
