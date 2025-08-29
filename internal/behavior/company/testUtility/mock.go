package testUtility

import (
	entryAuth "github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type UserCreatorMock struct {
	FakeExecute func(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error)
}

func (mock UserCreatorMock) Execute(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error) {
	return mock.FakeExecute(entry)
}

type UserGetterMock struct {
	FakeExecute func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error)
}

func (mock UserGetterMock) Execute(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
	return mock.FakeExecute(identifier)
}

type EmailSetterMock struct {
	FakeExecute func(entry entryUser.EmailGetter, user *shelterUser.UserAuthentic) error
}

func (mock EmailSetterMock) Execute(entry entryUser.EmailGetter, user *shelterUser.UserAuthentic) error {
	return mock.FakeExecute(entry, user)
}

type EmailVerifierMock struct {
	FakeExecute func(entry entryUser.EmailVerifier, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error)
}

func (mock EmailVerifierMock) Execute(entry entryUser.EmailVerifier, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error) {
	return mock.FakeExecute(entry, user)
}

type NameChangerMock struct {
	FakeExecute func(entry entryUser.UserApplyer, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error)
}

func (mock NameChangerMock) Execute(entry entryUser.UserApplyer, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error) {
	return mock.FakeExecute(entry, user)
}

type PasswordSetterMock struct {
	FakeExecute func(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error
}

func (mock PasswordSetterMock) Execute(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error {
	return mock.FakeExecute(entry, user)
}

type PasswordCheckerMock struct {
	FakeExecute func(entry entryAuth.AuthLoginner) (pkgText.Identifier, error)
}

func (mock PasswordCheckerMock) Execute(entry entryAuth.AuthLoginner) (pkgText.Identifier, error) {
	return mock.FakeExecute(entry)
}

type RefreshTokenIssuerMock struct {
	FakeExecute func(user *shelterUser.UserAuthentic) (shelterText.Token, error)
}

func (mock RefreshTokenIssuerMock) Execute(user *shelterUser.UserAuthentic) (shelterText.Token, error) {
	return mock.FakeExecute(user)
}

type RefreshTokenCheckerMock struct {
	FakeExecute func(entry entryAuth.RefreshTokenGetter) (*shelterUser.UserAuthentic, error)
}

func (mock RefreshTokenCheckerMock) Execute(entry entryAuth.RefreshTokenGetter) (*shelterUser.UserAuthentic, error) {
	return mock.FakeExecute(entry)
}

type AccessTokenIssuerMock struct {
	FakeExecute func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error)
}

func (mock AccessTokenIssuerMock) Execute(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
	return mock.FakeExecute(user)
}
