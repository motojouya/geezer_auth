package accessToken

import (
	"time"
)

// model/userとの変換が必要。
// accessToken->modelは、あんまりなさそう。というか、exposeIdからdbのデータ引っ張ってくるほうが現実的
// model->accessTokenは、accessTokenを作る際に必要なので、実装がほしい。
// 変換ロジックとしてどこかに持つか、あるいはprocedureの中に持たせてもいい。
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
