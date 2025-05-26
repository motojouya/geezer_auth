package user

type Company struct {
	ExposeId ExposeId
	Name     Name
}

func NewCompany(exposeId string, name string) Company {
	return Company{
		ExposeId: exposeId,
		Name: name,
	}
}
