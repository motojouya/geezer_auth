package command_test

import (
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	"github.com/motojouya/geezer_auth/internal/db/query/command"
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExpireEmailExcludeVerified(t *testing.T) {
	var now = testUtility.GetNow()
	testUtility.Truncate(t, orp)

	// user tableはforeign key制約があるのでいれとく
	var userRecords = []user.User{
		//           persist_key, identifier , email_idetifier     , name       , bot_flag  , register_date        , update_date
		user.NewUser(0 /*     */, "US-TASTAS", "test01@example.com", "tast name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TESTES", "test02@example.com", "test name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TOSTOS", "test03@example.com", "tost name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
	}
	var savedUserRecords = testUtility.Ready(t, orp, userRecords)

	var pastVerifyDate = now.AddDate(0, 0, -3)
	var pastExpireDate = now.AddDate(0, 0, -7)

	var records = []*user.UserEmail{
		//                persist_key, user_persist_key              , email               , verify_token01  ,register_date,verify_date, expire_date
		user.NewUserEmail(0 /*     */, savedUserRecords[0].PersistKey, "test01@example.com", "verify_token01", now /*    */, nil /*       */, nil),             // x user不一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test02@example.com", "verify_token02", now /*    */, nil /*       */, nil),             // x persist_key一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test03@example.com", "verify_token03", now /*    */, nil /*       */, nil),             // o verify_date null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test04@example.com", "verify_token04", now /*    */, &pastVerifyDate, nil),             // x verify_date not null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test05@example.com", "verify_token05", now /*    */, nil /*       */, &pastExpireDate), // x expired
	}
	var savedEmailRecords = testUtility.ReadyPointer(t, orp, records)

	var err = command.ExpireEmail(orp, savedUserRecords[1].PersistKey, savedEmailRecords[1].PersistKey, true, now)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	var expectRecords = []*user.UserEmail{
		//                persist_key, user_persist_key              , email               , verify_token01  ,register_date,verify_date, expire_date
		user.NewUserEmail(0 /*     */, savedUserRecords[0].PersistKey, "test01@example.com", "verify_token01", now /*    */, nil /*       */, nil),             // x user不一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test02@example.com", "verify_token02", now /*    */, nil /*       */, nil),             // x persist_key一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test03@example.com", "verify_token03", now /*    */, nil /*       */, &now),            // o verify_date null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test04@example.com", "verify_token04", now /*    */, &pastVerifyDate, nil),             // x verify_date not null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test05@example.com", "verify_token05", now /*    */, nil /*       */, &pastExpireDate), // x expired
	}

	var expects = essence.ToVal(expectRecords)
	testUtility.AssertTable(t, orp, []string{"persist_key"}, expects, assertSameUserEmail)
}

func TestExpireEmailIncludeVerified(t *testing.T) {
	var now = testUtility.GetNow()
	testUtility.Truncate(t, orp)

	// user tableはforeign key制約があるのでいれとく
	var userRecords = []user.User{
		//           persist_key, identifier , email_idetifier     , name       , bot_flag  , register_date        , update_date
		user.NewUser(0 /*     */, "US-TASTAS", "test01@example.com", "tast name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TESTES", "test02@example.com", "test name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TOSTOS", "test03@example.com", "tost name", false /**/, now.AddDate(0, -1, 0), now.AddDate(0, 0, -3)),
	}
	var savedUserRecords = testUtility.Ready(t, orp, userRecords)

	var pastVerifyDate = now.AddDate(0, 0, -3)
	var pastExpireDate = now.AddDate(0, 0, -7)

	var records = []*user.UserEmail{
		//                persist_key, user_persist_key              , email               , verify_token01  ,register_date,verify_date, expire_date
		user.NewUserEmail(0 /*     */, savedUserRecords[0].PersistKey, "test01@example.com", "verify_token01", now /*    */, nil /*       */, nil),             // x user不一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test02@example.com", "verify_token02", now /*    */, nil /*       */, nil),             // x persist_key一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test03@example.com", "verify_token03", now /*    */, nil /*       */, nil),             // o verify_date null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test04@example.com", "verify_token04", now /*    */, &pastVerifyDate, nil),             // o verify_date not null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test05@example.com", "verify_token05", now /*    */, nil /*       */, &pastExpireDate), // x expired
	}
	var savedEmailRecords = testUtility.ReadyPointer(t, orp, records)

	var err = command.ExpireEmail(orp, savedUserRecords[1].PersistKey, savedEmailRecords[1].PersistKey, false, now)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	var expectRecords = []*user.UserEmail{
		//                persist_key, user_persist_key              , email               , verify_token01  ,register_date,verify_date, expire_date
		user.NewUserEmail(0 /*     */, savedUserRecords[0].PersistKey, "test01@example.com", "verify_token01", now /*    */, nil /*       */, nil),             // x user不一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test02@example.com", "verify_token02", now /*    */, nil /*       */, nil),             // x persist_key一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test03@example.com", "verify_token03", now /*    */, nil /*       */, &now),            // o verify_date null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test04@example.com", "verify_token04", now /*    */, &pastVerifyDate, &now),            // o verify_date not null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test05@example.com", "verify_token05", now /*    */, nil /*       */, &pastExpireDate), // x expired
	}

	var expects = essence.ToVal(expectRecords)
	testUtility.AssertTable(t, orp, []string{"persist_key"}, expects, assertSameUserEmail)
}

func assertSameUserEmail(t *testing.T, expect user.UserEmail, actual user.UserEmail) {
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
}
