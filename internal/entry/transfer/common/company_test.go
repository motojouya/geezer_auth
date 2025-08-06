package common_test

import (
	"github.com/motojouya/geezer_auth/internal/entry/transfer/common"
	"github.com/motojouya/geezer_auth/internal/shelter/company"
	"github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFromShelterCompany(t *testing.T) {
	var persistKey uint = 1
	var identifier, _ = text.NewIdentifier("CP-TESTES")
	var name, _ = text.NewName("TestRole")
	var registeredDate = time.Now()

	var shelterCompany = company.NewCompany(persistKey, identifier, name, registeredDate)

	var transferCompany = common.FromShelterCompany(shelterCompany)

	assert.Equal(t, string(identifier), transferCompany.Identifier)
	assert.Equal(t, string(name), transferCompany.Name)

	t.Logf("transferCompany: %+v", transferCompany)
}
