package accessToken

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"os"
	"strconv"
)

// TODO middlewareも作ってしまいたい。イメージを掴んで置く
// TODO 本packageは、externalにするので、どこからも依存しない感じにしておく。
// 変換としては、model.Userから、このuserへの変換が必要だが、それはmodelに生やす感じで

type Company struct {
	ExposeId string
	Name string
	Role string
}

func NewCompany(exposeId string, name string, role string) *Company {
	return &Company{
		ExposeId: exposeId,
		Name: name,
		Role: role,
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
