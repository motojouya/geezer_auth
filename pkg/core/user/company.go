package user

type Company struct {
	ExposeId ExposeId
	Name     Name
}

func NewCompany(exposeId ExposeId, name Name) Company {
	return Company{
		ExposeId: exposeId,
		Name:     name,
	}
}
