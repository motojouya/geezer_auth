package company_test

import (
	//"errors"
	"github.com/motojouya/geezer_auth/internal/behavior/company"
	dbCompany "github.com/motojouya/geezer_auth/internal/db/transfer/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type companyGetterDBMock struct {
	getCompany func(identifier string) (*dbCompany.Company, error)
}

func (mock companyGetterDBMock) GetCompany(identifier string) (*dbCompany.Company, error) {
	return mock.getCompany(identifier)
}

func getDbCompanyForGet(expectId string) dbCompany.Company {
	return dbCompany.Company{
		PersistKey:     1,
		Identifier:     expectId,
		Name:           "TestCompany",
		RegisteredDate: time.Now(),
	}
}

func getCompanyGetterDbMock(t *testing.T, expectId string) companyGetterDBMock {
	var dbCompanyVal = getDbCompanyForGet(expectId)
	var getCompany = func(identifier string) (*dbCompany.Company, error) {
		assert.Equal(t, expectId, identifier)
		return &dbCompanyVal, nil
	}
	return companyGetterDBMock{
		getCompany: getCompany,
	}
}

type companyGetterEntryMock struct {
	getCompanyIdentifier func() (pkgText.Identifier, error)
}

func (mock companyGetterEntryMock) GetCompanyIdentifier() (pkgText.Identifier, error) {
	return mock.getCompanyIdentifier()
}

func getCompanyGetterEntryMock(t *testing.T, expectId string) companyGetterEntryMock {
	var getCompanyIdentifier = func() (pkgText.Identifier, error) {
		return pkgText.Identifier(expectId), nil
	}
	return companyGetterEntryMock{
		getCompanyIdentifier: getCompanyIdentifier,
	}
}

func TestCompanyGetter(t *testing.T) {
	var expectId = "CP-TESTES"

	var dbMock = getCompanyGetterDbMock(t, expectId)
	var entryMock = getCompanyGetterEntryMock(t, expectId)

	getter := company.NewCompanyGet(dbMock)
	result, err := getter.Execute(entryMock)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectId, string(result.Identifier))

	t.Logf("get company: %+v", result)
}

// TODO working error cases
