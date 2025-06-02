package authorization

type RequirePermission struct {
	SelfEdit      bool
	CompanyAccess bool
	CompanyInvite bool
	CompanyEdit   bool
}

func NewRequirePermission(selfEdit bool, companyAccess bool, companyInvite bool, companyEdit bool) RequirePermission {
	return RequirePermission{
		SelfEdit:      selfEdit,
		CompanyAccess: companyAccess,
		CompanyInvite: companyInvite,
		CompanyEdit:   companyEdit,
	}
}
