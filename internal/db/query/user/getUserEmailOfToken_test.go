package user_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserEmailOfToken(t *testing.T) {
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
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test01@example.com", "verify_token01", now /*    */, nil /*  */, &now), // o 一致 expireは関係ない
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test02@example.com", "verify_token01", now /*    */, nil /*  */, nil),  // x email不一致
		user.NewUserEmail(0 /*     */, savedUserRecords[0].PersistKey, "test01@example.com", "verify_token01", now /*    */, nil /*  */, nil),  // x user不一致
	}
	testUtility.ReadyPointer(t, orp, records)

	var result, err = orp.GetUserEmailOfToken("US-TESTES", "test01@example.com")
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	var expect = user.UserEmailFull{
		UserEmail: user.UserEmail{
			PersistKey:     1,
			UserPersistKey: savedUserRecords[1].PersistKey,
			Email:          "test01@example.com",
			VerifyToken:    "verify_token01",
			RegisterDate:   now,
			VerifyDate:     nil,
			ExpireDate:     &now,
		},
		UserIdentifier:     "US-TESTES",
		UserExposeEmailId:  "test02@example.com",
		UserName:           "test name",
		UserBotFlag:        false,
		UserRegisteredDate: now.AddDate(0, -1, 0),
		UserUpdateDate:     now.AddDate(0, 0, -3),
	}

	assert.NotNil(t, result)
	assertSameUserEmail(t, expect, *result)
}
