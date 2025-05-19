package accessToken

import (
	"time"
)

type Company struct {
	ExposeId string
	Name     string
	Role     string
	RoleName string
}

func NewCompany(exposeId string, name string, role string, roleName string) *Company {
	return &Company{
		ExposeId: exposeId,
		Name: name,
		Role: role,
		RoleName: roleName,
	}
}

type User struct {
	ExposeId string
	EmailId string
	Email *string
	Name string
	BotFlag bool
	Company *Company
	UpdateDate time.Time
}

func NewUser(exposeId string, emailId string, email *string, name string, botFlag bool, company *Company, updateDate time.Time) *User {
	return &User{
		ExposeId: exposeId,
		EmailId: emailId,
		Email: email,
		Name: name,
		BotFlag: botFlag,
		Company: company,
		UpdateDate: updateDate,
	}
}
