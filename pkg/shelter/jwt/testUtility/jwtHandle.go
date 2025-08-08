package testUtility

import (
	"time"
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/motojouya/geezer_auth/pkg/shelter/user"
)

type JwtHandlerMock struct {
	FakeGenerate func(user *user.User, issueDate time.Time, id string) (*user.Authentic, text.JwtToken, error)
}
	
func (mock JwtHandlerMock) Generate(user *user.User, issueDate time.Time, id string) (*user.Authentic, text.JwtToken, error) {
	return mock.FakeGenerate(user, issueDate, id)
}
