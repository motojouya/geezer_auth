package user_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetUserEmail(t *testing.T) {
	var now = testUtility.GetNow()
	testUtility.Truncate(t, orp)

	var userRecords = []user.User{
		//           persist_key, identifier , email_idetifier     , name       , bot_flag  , register_date        , update_date
		user.NewUser(0 /*     */, "US-TASTAS", "test01@example.com", "tast name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TESTES", "test02@example.com", "test name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TOSTOS", "test03@example.com", "tost name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
	}
	var savedUserRecords = testUtility.Ready(t, orp, userRecords)

	var records = []*user.UserEmail{
		//                persist_key, user_persist_key              , email               , verify_token01  ,register_date,verify_date, expire_date
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test01@example.com", "verify_token01", now /*    */, nil /*  */, nil), //  o email一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test02@example.com", "verify_token01", now /*    */, nil /*  */, nil), //  x email不一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test01@example.com", "verify_token01", now /*    */, nil /*  */, &now), // x expired
	}
	testUtility.ReadyPointer(t, orp, records)

	var result, err = orp.GetUserEmail("test01@example.com")
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	var expects = []user.UserEmailFull{
		user.UserEmailFull{
			UserEmail: user.UserEmail{
				PersistKey:     1,
				UserPersistKey: savedUserRecords[1].PersistKey,
				Email:          "test01@example.com",
				VerifyToken:    "verify_token01",
				RegisterDate:   now,
				VerifyDate:     nil,
				ExpireDate:     nil,
			},
			UserIdentifier:        "US-TESTES",
			UserExposeEmailId:     "test02@example.com",
			UserName:              "test name",
			UserBotFlag:           false,
			UserRegisteredDate:    now.AddDate(0, -1, 0),
			UserUpdateDate:        now.AddDate(0, 0, -3),
		},
	}

	testUtility.AssertRecords(t, expects, result, assertSameUserEmail)
}

func assertSameUserEmail(t *testing.T, expect user.UserEmailFull, actual user.UserEmailFull) {
	assert.Equal(t, expect.UserPersistKey, actual.UserPersistKey)
	assert.Equal(t, expect.Email, actual.Email)
	assert.Equal(t, expect.VerifyToken, actual.VerifyToken)
	assert.WithinDuration(t, expect.RegisterDate, actual.RegisterDate, time.Second)
	if expect.VerifyDate == nil {
		assert.Nil(t, actual.VerifyDate)
	} else {
		assert.WithinDuration(t, *expect.VerifyDate, *actual.VerifyDate, time.Second)
	}
	if expect.ExpireDate == nil {
		assert.Nil(t, actual.ExpireDate)
	} else {
		assert.WithinDuration(t, *expect.ExpireDate, *actual.ExpireDate, time.Second)
	}

	assert.Equal(t, expect.UserIdentifier, actual.UserIdentifier)
	assert.Equal(t, expect.UserExposeEmailId, actual.UserExposeEmailId)
	assert.Equal(t, expect.UserName, actual.UserName)
	assert.Equal(t, expect.UserBotFlag, actual.UserBotFlag)
	assert.WithinDuration(t, expect.UserRegisteredDate, actual.UserRegisteredDate, time.Second)
	assert.WithinDuration(t, expect.UserUpdateDate, actual.UserUpdateDate, time.Second)
}
