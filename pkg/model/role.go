package model

type Role {
	Label Label
	Name  Name
}

func NewRole(label Label, name Name) Role {
	return Role{
		label:   label,
		Name:    name,
	}
}
