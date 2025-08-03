package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/motojouya/geezer_auth/pkg/shelter/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCompany(t *testing.T) {
	var identifier, _ = text.NewIdentifier("CP-TESTES")
	var name, _ = text.NewName("TestCompany")

	var company = user.NewCompany(identifier, name)

	assert.Equal(t, string(name), string(company.Name))
	assert.Equal(t, string(identifier), string(company.Identifier))

	t.Logf("company: %+v", company)
	t.Logf("company.Identifier: %s", string(company.Identifier))
	t.Logf("company.Name: %s", string(company.Name))
}
