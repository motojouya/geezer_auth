package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewCompany(t *testing.T) {
	var exposeId, _ = text.NewExposeId("CP-TESTES")
	var name, _ = text.NewExposeId("TestCompany")

	var company = user.NewCompany(exposeId, name, role, roleName)

	assert.Equal(t, string(name), string(company.Name))
	assert.Equal(t, string(exposeId), string(company.ExposeId))

	t.Logf("company: %+v", company)
	t.Logf("company.ExposeId: %s", string(company.ExposeId))
	t.Logf("company.Name: %s", string(company.Name))
}
