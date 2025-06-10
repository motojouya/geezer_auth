package user

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
)

type Role struct {
	Label text.Label
	Name  text.Name
}

func NewRole(label text.Label, name text.Name) Role {
	return Role{
		label: label,
		Name:  name,
	}
}
