package common_test

import (
	"github.com/motojouya/geezer_auth/internal/core/company"
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	"github.com/motojouya/geezer_auth/pkg/core/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromCoreCompany(t *testing.T) {
	var persistKey uint = 1
	var identifier, _ = text.NewIdentifier("CP-TESTES")
	var name, _ = text.NewName("TestRole")
	var registeredDate = time.Now()

	var coreCompany = company.NewCompany(persistKey, identifier, name, registeredDate)

	var transferCompany = common.FromCoreCompany(coreCompany)

	assert.Equal(t, string(identifier), transferCompany.Identifier)
	assert.Equal(t, string(name), transferCompany.Name)

	t.Logf("transferCompany: %+v", transferCompany)
}
