package user_test

import (
	"github.com/google/uuid"
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	shelterText "github.com/motojouya/geezer_auth/internal/shelter/text"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	dbUser "github.com/motojouya/geezer_auth/internal/db/transfer/user"
	entryUser "github.com/motojouya/geezer_auth/internal/entry/transfer/user"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	pkgText "github.com/motojouya/geezer_auth/pkg/shelter/text"
	"github.com/motojouya/geezer_auth/internal/behavior/user"
	"github.com/motojouya/geezer_auth/internal/behavior/testUtility"
)

type UserCreatorDB struct {
	testUtility.SqlExecutorMock
	t *testing.T
	expectIdentifier string
	resultUser *dbUser.User
	resultAuthentic *dbUser.UserAuthentic
}

func (db UserCreatorDB) GetUser(identifier string) (*dbUser.User, error) {
	if db.expectIdentifier != "" && identifier != db.expectIdentifier {
		db.t.Errorf("Expected identifier %s, got %s", db.expectIdentifier, identifier)
		return nil, nil
	}

	return db.resultUser, nil
}

func (db UserCreatorDB) GetUserAuthentic(identifier string, now time.Time) (*dbUser.UserAuthentic, error) {
	if db.expectIdentifier != "" && identifier != db.expectIdentifier {
		db.t.Errorf("Expected identifier %s, got %s", db.expectIdentifier, identifier)
		return nil, nil
	}

	return db.resultAuthentic, nil
}

type userEntry struct {
}

func (ue userEntry) ToCoreUser(identifier pkgText.Identifier, now time.Time) (*shelterUser.UnsavedUser, error) {
	return &shelterUser.UnsavedUser{
		Identifier: identifier,
		Name:       "Test User",
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func TestUserCreate(t *testing.T) {
	now := time.Now()
	localMock := testUtility.NewLocalerMock(t, "TESTES", pkgText.IdentifierLength, pkgText.IdentifierChar, uuid.UUID("UUIDTESTuuidtestUUIDTESTuuidtest"), nil, now)
	resultUser := &dbUser.User{
		Identifier: "US-TESTES",
		Name: "Test User",
		CreatedAt: now,
		UpdatedAt: now,
	}
	resultAuthentic := &dbUser.UserAuthentic{
		Identifier: "US-TESTES",
		Name: "Test User",
		CreatedAt: now,
		UpdatedAt: now,
	}
	dbMock := UserCreatorDB{
		SqlExecutorMock: testUtility.SqlExecutorMock{},
		t: t,
		expectIdentifier: "US-TESTES",
		resultUser: resultUser,
		resultAuthentic: resultAuthentic,
	}

	creator := user.NewUserCreate(localMock, dbMock)
	result := creator.Execute(userEntry{})

	// Mock the local time
	local.On("GetNow").Return(time.Now())

	// Mock the user identifier creation
	local.On("GenerateRamdomString", pkgText.IdentifierLength, pkgText.IdentifierChar).Return("randomString")

	// Mock the database interactions
	db.On("GetUser", "randomString").Return(nil, nil)

	// Create a user entry
	entry := &entryUser.UserGetter{
		Name: "Test User",
	}

	// Execute the user creation
	userAuthentic, err := creator.Execute(entry)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, userAuthentic)
	assert.Equal(t, "Test User", userAuthentic.Name)
}

