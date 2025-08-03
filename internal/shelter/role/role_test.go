package role_test

import (
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	pkg "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateRole(t *testing.T) {
	var name, _ = pkg.NewName("TestRole")
	var label, _ = pkg.NewLabel("TEST_ROLE")
	var description, _ = text.NewText("Role for testing")
	var registeredDate = time.Now()

	var role = role.NewRole(name, label, description, registeredDate)

	assert.Equal(t, string(name), string(role.Name))
	assert.Equal(t, string(label), string(role.Label))
	assert.Equal(t, string(description), string(role.Description))
	assert.Equal(t, registeredDate, role.RegisteredDate)

	t.Logf("role: %+v", role)
	t.Logf("role.Name: %s", role.Name)
	t.Logf("role.Label: %s", role.Label)
	t.Logf("role.Description: %s", role.Description)
	t.Logf("role.RegisteredDate: %s", role.RegisteredDate)
}
