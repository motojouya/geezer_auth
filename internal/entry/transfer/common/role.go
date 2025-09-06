package common

import (
	shelter "github.com/motojouya/geezer_auth/internal/shelter/role"
)

type Role struct {
	Label       string `json:"label"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func FromShelterRole(r shelter.Role) Role {
	return Role{
		Label:       string(r.Label),
		Name:        string(r.Name),
		Description: string(r.Description),
	}
}
