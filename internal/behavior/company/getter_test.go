package company_test

import (
	//"errors"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type companyGetterDBMock struct {
	getCompany func(identifier string) (*dbCompany.Company, error)
}

func (mock companyGetterDBMock) getCompany(identifier string) (*dbCompany.Company, error) {
	return mock.getCompany(identifier)
}

func getDbCompanyForGet(expectId string) dbCompany.Company {
	return dbCompany.Company{
		Id:             1,
		Identifier:     expectId,
		Name:           "TestCompany",
		RegisteredDate: registeredDate,
	}
}

func getCompanyGetterDbMock(t *testing.T, expectId string) companyGetterDBMock {
	var dbCompany = getDbCompanyForGet(expectId)
	var getCompany = func(identifier string) (*dbCompany.Company, error) {
		assert.Equal(t, expectId, identifier)
		return dbCompany, nil
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

func TestUserGetter(t *testing.T) {
	var expectId = "CP-TESTES"

	var dbMock = getCompanyGetterDbMock(t, expectId)
	var entryMock = getCompanyGetterEntryMock(t, expectId)

	getter := company.NewCompanyGet(localerMock, dbMock)
	result, err := getter.Execute(entryMock)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectId, string(result.Identifier))

	t.Logf("get company: %+v", result)
}

// TODO working error cases
