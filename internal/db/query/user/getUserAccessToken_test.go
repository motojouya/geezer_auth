package user_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetUserAccessToken(t *testing.T) {
	var now = testUtility.GetNow()
	testUtility.Truncate(t, orp)

	var userRecords = []user.User{
		//           persist_key, identifier , email_idetifier     , name       , bot_flag  , register_date        , update_date
		user.NewUser(0 /*     */, "US-TASTAS", "test01@example.com", "tast name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TESTES", "test01@example.com", "test name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TOSTOS", "test01@example.com", "tost name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
	}
	var savedUserRecords = testUtility.Ready(t, orp, userRecords)

	var records = []user.UserAccessToken{
		//                      persist_key, user_persist_key              , access_token    , source_update_date   , register_date        , expire_date
		user.NewUserAccessToken(0 /*     */, savedUserRecords[0].PersistKey, "access_token05", now.AddDate(0, 0, -3), now.AddDate(0, 0, -3), now.AddDate(0, 0, 7)),  // x user.identifierと一致しない
		user.NewUserAccessToken(0 /*     */, savedUserRecords[1].PersistKey, "access_token04", now.AddDate(0, 0, -3), now.AddDate(0, 0, -4), now.AddDate(0, 0, 7)),  // o 取得対象
		user.NewUserAccessToken(0 /*     */, savedUserRecords[1].PersistKey, "access_token03", now.AddDate(0, 0, -3), now.AddDate(0, 0, -5), now.AddDate(0, 0, 7)),  // o 取得対象
		user.NewUserAccessToken(0 /*     */, savedUserRecords[1].PersistKey, "access_token02", now.AddDate(0, 0, -7), now.AddDate(0, 0, -6), now.AddDate(0, 0, 7)),  // x user.update_dateと一致しない
		user.NewUserAccessToken(0 /*     */, savedUserRecords[1].PersistKey, "access_token01", now.AddDate(0, 0, -3), now.AddDate(0, 0, -7), now.AddDate(0, 0, -7)), // x expire_dateが過去
	}
	testUtility.Ready(t, orp, records)

	var result, err = orp.GetUserAccessToken("US-TESTES", now)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	var expects = []user.UserAccessTokenFull{
		user.UserAccessTokenFull{
			UserAccessToken: user.UserAccessToken{
				PersistKey:       1,
				UserPersistKey:   savedUserRecords[1].PersistKey,
				AccessToken:      "access_token04",
				SourceUpdateDate: now.AddDate(0, 0, -3),
				RegisterDate:     now.AddDate(0, 0, -4),
				ExpireDate:       now.AddDate(0, 0, 7),
			},
			UserIdentifier:     "US-TESTES",
			UserExposeEmailId:  "test01@example.com",
			UserName:           "test name",
			UserBotFlag:        false,
			UserRegisteredDate: now.AddDate(0, -1, 0),
			UserUpdateDate:     now.AddDate(0, 0, -3),
		},
		user.UserAccessTokenFull{
			UserAccessToken: user.UserAccessToken{
				PersistKey:       2,
				UserPersistKey:   savedUserRecords[1].PersistKey,
				AccessToken:      "access_token03",
				SourceUpdateDate: now.AddDate(0, 0, -3),
				RegisterDate:     now.AddDate(0, 0, -5),
				ExpireDate:       now.AddDate(0, 0, 7),
			},
			UserIdentifier:     "US-TESTES",
			UserExposeEmailId:  "test01@example.com",
			UserName:           "test name",
			UserBotFlag:        false,
			UserRegisteredDate: now.AddDate(0, -1, 0),
			UserUpdateDate:     now.AddDate(0, 0, -3),
		},
	}

	testUtility.AssertRecords(t, expects, result, assertSameUserAccessToken)
}

func assertSameUserAccessToken(t *testing.T, expect user.UserAccessTokenFull, actual user.UserAccessTokenFull) {
	assert.Equal(t, expect.UserPersistKey, actual.UserPersistKey)
	assert.Equal(t, expect.UserIdentifier, actual.UserIdentifier)
	assert.Equal(t, expect.AccessToken, actual.AccessToken)
	assert.WithinDuration(t, expect.SourceUpdateDate, actual.SourceUpdateDate, time.Second)
	assert.WithinDuration(t, expect.RegisterDate, actual.RegisterDate, time.Second)
	assert.WithinDuration(t, expect.ExpireDate, actual.ExpireDate, time.Second)
	assert.Equal(t, expect.UserIdentifier, actual.UserIdentifier)
	assert.Equal(t, expect.UserExposeEmailId, actual.UserExposeEmailId)
	assert.Equal(t, expect.UserName, actual.UserName)
	assert.Equal(t, expect.UserBotFlag, actual.UserBotFlag)
	assert.WithinDuration(t, expect.UserRegisteredDate, actual.UserRegisteredDate, time.Second)
	assert.WithinDuration(t, expect.UserUpdateDate, actual.UserUpdateDate, time.Second)
}
