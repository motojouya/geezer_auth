package user_test

import (
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


func getLocalerMock(t *testing.T, randomString string, expectLength int, expectSource string, uuid uuid.UUID, uuidErr error, now time.Time) localPkg.Localer {
	return testUtility.NewLocalerMock(t, randomString, expectLength, expectSource, uuid, uuidErr, now)
}

func TestUserCreate(t *testing.T) {
	// Mock dependencies
	local := &localPkg.MockLocaler{}
	db := &userQuery.MockUserCreatorDB{}

	// Create a UserCreate instance
	creator := user.NewUserCreate(local, db)

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

