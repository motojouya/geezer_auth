package common_test

import (
	"github.com/motojouya/geezer_auth/internal/core/role"
	"github.com/motojouya/geezer_auth/internal/core/text"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	pkg "github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromCoreRole(t *testing.T) {
	var name, _ = pkg.NewName("TestRole")
	var label, _ = pkg.NewLabel("TEST_ROLE")
	var description, _ = text.NewText("Role for testing")
	var registeredDate = time.Now()

	var coreRole = role.NewRole(name, label, description, registeredDate)

	var transferRole = common.FromCoreRole(coreRole)

	assert.Equal(t, string(name), transferRole.Name)
	assert.Equal(t, string(label), transferRole.Label)
	assert.Equal(t, string(description), transferRole.Description)

	t.Logf("transferRole: %+v", transferRole)
}
