package user_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetUserRefreshToken(t *testing.T) {
	var now = testUtility.GetNow()
	testUtility.Truncate(t, orp)

	var userRecords = []user.User{
		//           persist_key, identifier , email_idetifier     , name       , bot_flag  , register_date        , update_date
		user.NewUser(0 /*     */, "US-TASTAS", "test01@example.com", "tast name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TESTES", "test02@example.com", "test name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TOSTOS", "test03@example.com", "tost name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
	}
	var savedUserRecords = testUtility.Ready(t, orp, userRecords)

	var records = []user.UserRefreshToken{
		//                       persist_key, user_persist_key              , refresh_token    , register_date        , expire_date
		user.NewUserRefreshToken(0 /*     */, savedUserRecords[0].PersistKey, "refresh_token01", now.AddDate(0, -1, 0), now.AddDate(0, 0, 7)),  // x user不一致
		user.NewUserRefreshToken(0 /*     */, savedUserRecords[1].PersistKey, "refresh_token02", now.AddDate(0, -1, 0), now.AddDate(0, 0, 7)),  // o 対象
		user.NewUserRefreshToken(0 /*     */, savedUserRecords[1].PersistKey, "refresh_token03", now.AddDate(0, -1, 0), now.AddDate(0, 0, -7)), // x expire_date過去
	}
	testUtility.Ready(t, orp, records)

	var result, err = orp.GetUserRefreshToken("US-TESTES", now)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	var expect = user.UserRefreshTokenFull{
		UserRefreshToken: user.UserRefreshToken{
			PersistKey:     1,
			UserPersistKey: savedUserRecords[1].PersistKey,
			RefreshToken:   "refresh_token02",
			RegisterDate:   now.AddDate(0, -1, 0),
			ExpireDate:     now.AddDate(0, 0, 7),
		},
		UserIdentifier:     "US-TESTES",
		UserExposeEmailId:  "test02@example.com",
		UserName:           "test name",
		UserBotFlag:        false,
		UserRegisteredDate: now.AddDate(0, -1, 0),
		UserUpdateDate:     now.AddDate(0, 0, -3),
	}

	assert.NotNil(t, result)
	assertSameUserRefreshToken(t, expect, *result)
}

func assertSameUserRefreshToken(t *testing.T, expect user.UserRefreshTokenFull, actual user.UserRefreshTokenFull) {
	assert.Equal(t, expect.UserPersistKey, actual.UserPersistKey)
	assert.Equal(t, expect.RefreshToken, actual.RefreshToken)
	assert.WithinDuration(t, expect.RegisterDate, actual.RegisterDate, time.Second)
	assert.WithinDuration(t, expect.ExpireDate, actual.ExpireDate, time.Second)

	assert.Equal(t, expect.UserIdentifier, actual.UserIdentifier)
	assert.Equal(t, expect.UserExposeEmailId, actual.UserExposeEmailId)
	assert.Equal(t, expect.UserName, actual.UserName)
	assert.Equal(t, expect.UserBotFlag, actual.UserBotFlag)
	assert.WithinDuration(t, expect.UserRegisteredDate, actual.UserRegisteredDate, time.Second)
	assert.WithinDuration(t, expect.UserUpdateDate, actual.UserUpdateDate, time.Second)
}
