package model

import (
	"time"
)

type UnsavedRole struct {
	Name        string
	Label       string
	Description string
}

type Role struct {
	RoleId         uint
	RegisteredDate time.Time
	UnsavedRole
}

func CreateRole(name string, label string, description string) UnsavedRole {
	return UnsavedRole{
		Name:        name,
		Label:       label,
		Description: description,
	}
}

func NewRole(roleId uint, name string, label string, description string, registeredDate time.Time) Role {
	return Role{
		RoleId:         roleId,
		RegisteredDate: registeredDate,
		UnsavedRole:    CreateRole(name, label, description),
	}
}
