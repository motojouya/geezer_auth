package transfer

import (
	core "github.com/motojouya/geezer_auth/internal/core/role"
)

type Role struct {
	Label          string
	Name           string
	Description    string
}

func FromCoreRole(r core.Role) Role {
	return Role{
		Label:       string(r.Label),
		Name:        string(r.Name),
		Description: string(r.Description),
	}
}
