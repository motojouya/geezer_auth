package user_test

import (
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewRole(t *testing.T) {
	var name, _ = text.NewName("TestCompany")
	var label, _ = text.NewLabel("TestLabel")

	var role = user.CreateRole(label, name)

	assert.Equal(t, string(name), string(role.Name))
	assert.Equal(t, string(label), string(role.Label))

	t.Logf("company: %+v", role)
	t.Logf("company.Name: %s", string(role.Name))
	t.Logf("company.Label: %s", string(role.Label))
}
