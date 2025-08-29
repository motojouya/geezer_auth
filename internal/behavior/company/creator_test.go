package company_test

import (
	//"errors"
	"github.com/motojouya/geezer_auth/internal/behavior/company"
	dbUtility "github.com/motojouya/geezer_auth/internal/db/testUtility"
	dbCompany "github.com/motojouya/geezer_auth/internal/db/transfer/company"
	localUtility "github.com/motojouya/geezer_auth/internal/local/testUtility"
	shelterCompany "github.com/motojouya/geezer_auth/internal/shelter/company"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type companyCreatorDBMock struct {
	getCompany func(identifier string) (*dbCompany.Company, error)
	dbUtility.SqlExecutorMock
}

func (mock companyCreatorDBMock) GetCompany(identifier string) (*dbCompany.Company, error) {
	return mock.getCompany(identifier)
}

type companyEntryMock struct {
	toShelterCompany func(identifier pkgText.Identifier, now time.Time) (shelterCompany.UnsavedCompany, error)
}

func (mock companyEntryMock) ToShelterCompany(identifier pkgText.Identifier, now time.Time) (shelterCompany.UnsavedCompany, error) {
	return mock.toShelterCompany(identifier, now)
}

func getShelterCompanyForCreate(expectId string, expectName string) shelterCompany.UnsavedCompany {
	var identifier, _ = pkgText.NewIdentifier(expectId)
	var name, _ = pkgText.NewName(expectName)
	var registeredDate = time.Now()

	return shelterCompany.CreateCompany(identifier, name, registeredDate)
}

func getLocalerMockForCreate(t *testing.T, expectId string, now time.Time) *localUtility.LocalerMock {
	var getNow = func() time.Time {
		return now
	}
	var generateRamdomString = func(length int, charSet string) string {
		return expectId
	}
	return &localUtility.LocalerMock{
		FakeGenerateRamdomString: generateRamdomString,
		FakeGetNow:               getNow,
	}
}

func getDbCompany(expectId string, expectName string) dbCompany.Company {
	return dbCompany.Company{
		PersistKey:     1,
		Identifier:     expectId,
		Name:           expectName,
		RegisteredDate: time.Now(),
	}
}

func getCreateDbMock(t *testing.T, expectId string, expectName string, firstNow time.Time) companyCreatorDBMock {
	var callCount = 0
	var dbCompanyValue = getDbCompany(expectId, expectName)
	var getCompany = func(identifier string) (*dbCompany.Company, error) {
		if callCount == 0 {
			callCount++
			return nil, nil
		}
		assert.Equal(t, expectId, identifier, "Expected identifier 'US-TESTES'")
		return &dbCompanyValue, nil
	}
	var insert = func(args ...interface{}) error {
		assert.Equal(t, 1, len(args), "Expected 1 argument")

		company, ok := args[0].(*dbCompany.Company)
		if !ok {
			t.Errorf("Expected first argument to be of type *dbCompany.Company, got %T", args[0])
		}

		assert.NotNil(t, company)
		assert.Equal(t, expectId, company.Identifier)

		return nil
	}
	return companyCreatorDBMock{
		SqlExecutorMock: dbUtility.SqlExecutorMock{
			FakeInsert: insert,
		},
		getCompany: getCompany,
	}
}

func getEntryMock(t *testing.T, expectId string, expectName string, firstNow time.Time) companyEntryMock {
	var shelterUseCompanyrVal = getShelterCompanyForCreate(expectId, expectName)
	var toShelterCompany = func(identifier pkgText.Identifier, now time.Time) (shelterCompany.UnsavedCompany, error) {
		assert.Equal(t, expectId, string(identifier))
		assert.WithinDuration(t, now, firstNow, time.Second, "Expected 'now' to be within 1 second of current time")
		return shelterUseCompanyrVal, nil
	}
	return companyEntryMock{
		toShelterCompany: toShelterCompany,
	}
}

func TestCompanyCreate(t *testing.T) {
	var expectId = "CP-TESTES"
	var expectName = "TestCompany"
	var firstNow = time.Now()

	var localerMock = getLocalerMockForCreate(t, "TESTES", firstNow)
	var dbMock = getCreateDbMock(t, expectId, expectName, firstNow)
	var entryMock = getEntryMock(t, expectId, expectName, firstNow)

	creator := company.NewCompanyCreate(localerMock, dbMock)
	result, err := creator.Execute(entryMock)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectId, string(result.Identifier))
	assert.Equal(t, expectName, string(result.Name))

	t.Logf("Company created: %+v", result)
}

// TODO working error cases
