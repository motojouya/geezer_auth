package testUtility

import (
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
)

type userCreatorMock struct {
	execute func(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error)
}

func (mock userCreatorMock) Execute(entry entryUser.UserGetter) (*shelterUser.UserAuthentic, error) {
	return mock.execute(entry)
}

type userGetterMock struct {
	execute func(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error)
}

func (mock userGetterMock) Execute(identifier pkgText.Identifier) (*shelterUser.UserAuthentic, error) {
	return mock.execute(identifier)
}

type emailSetterMock struct {
	execute func(entry entryUser.EmailGetter, user *shelterUser.UserAuthentic) error
}

func (mock emailSetterMock) Execute(entry entryUser.EmailGetter, user *shelterUser.UserAuthentic) error {
	return mock.execute(entry, user)
}

type emailVerifierMock struct {
	execute func(entry entryUser.EmailVerifier, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error)
}

func (mock emailVerifierMock) Execute(entry entryUser.EmailVerifier, user *shelterUser.UserAuthentic) (*shelterUser.UserAuthentic, error) {
	return mock.execute(entry, user)
}

type passwordSetterMock struct {
	execute func(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error
}

func (mock passwordSetterMock) Execute(entry entryUser.PasswordGetter, user *shelterUser.UserAuthentic) error {
	return mock.execute(entry, user)
}

type refreshTokenIssuerMock struct {
	execute func(user *shelterUser.UserAuthentic) (shelterText.Token, error)
}

func (mock refreshTokenIssuerMock) Execute(user *shelterUser.UserAuthentic) (shelterText.Token, error) {
	return mock.execute(user)
}

type accessTokenIssuerMock struct {
	execute func(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error)
}

func (mock accessTokenIssuerMock) Execute(user *shelterUser.UserAuthentic) (pkgText.JwtToken, error) {
	return mock.execute(user)
}
