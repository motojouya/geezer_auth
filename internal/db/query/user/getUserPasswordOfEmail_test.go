package user_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserPasswordOfEmail(t *testing.T) {
	var now = testUtility.GetNow()
	testUtility.Truncate(t, orp)

	var userRecords = []user.User{
		//           persist_key, identifier , email_idetifier     , name       , bot_flag  , register_date        , update_date
		user.NewUser(0 /*     */, "US-TASTAS", "test01@example.com", "tast name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TESTES", "test02@example.com", "test name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TOSTOS", "test03@example.com", "tost name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
	}
	var savedUserRecords = testUtility.Ready(t, orp, userRecords)

	var records = []*user.UserPassword{
		//                   persist_key, user_persist_key              , password    ,register_date,expire_date
		user.NewUserPassword(0 /*     */, savedUserRecords[0].PersistKey, "password01", now /*    */, nil),  // x user 不一致
		user.NewUserPassword(0 /*     */, savedUserRecords[1].PersistKey, "password02", now /*    */, nil),  // o 対象
		user.NewUserPassword(0 /*     */, savedUserRecords[1].PersistKey, "password03", now /*    */, &now), // x expire
	}
	testUtility.ReadyPointer(t, orp, records)

	var result, err = orp.GetUserPasswordOfEmail("test02@example.com")
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	var expect = user.UserPasswordFull{
		UserPassword: user.UserPassword{
			PersistKey:     1,
			UserPersistKey: savedUserRecords[1].PersistKey,
			Password:       "password02",
			RegisteredDate: now,
			ExpireDate:     nil,
		},
		UserIdentifier:     "US-TESTES",
		UserExposeEmailId:  "test02@example.com",
		UserName:           "test name",
		UserBotFlag:        false,
		UserRegisteredDate: now.AddDate(0, -1, 0),
		UserUpdateDate:     now.AddDate(0, 0, -3),
	}

	assert.NotNil(t, result)
	assertSameUserPassword(t, expect, *result)
}
