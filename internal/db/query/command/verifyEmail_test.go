package command_test

import (
	"github.com/motojouya/geezer_auth/internal/core/essence"
	"github.com/motojouya/geezer_auth/internal/db/query/command"
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyEmail(t *testing.T) {
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
		user.NewUserEmail(0 /*     */, savedUserRecords[0].PersistKey, "test01@example.com", "verify_token01", now /*    */, nil /*       */, nil),             // user不一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test02@example.com", "verify_token02", now /*    */, nil /*       */, nil),             // 更新対象
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test03@example.com", "verify_token03", now /*    */, nil /*       */, nil),             // verify_date null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test04@example.com", "verify_token04", now /*    */, &pastVerifyDate, nil),             // verify_date not null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test05@example.com", "verify_token05", now /*    */, nil /*       */, &pastExpireDate), // expired
	}
	var savedEmailRecords = testUtility.ReadyPointer(t, orp, records)

	var updateRecord = savedEmailRecords[1]
	updateRecord.VerifyDate = &now

	var resultRecord, err = command.VerifyEmail(orp, updateRecord, now)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	assert.NotNil(t, resultRecord)
	assertSameUserEmail(t, *updateRecord, *resultRecord)

	var expectRecords = []*user.UserEmail{
		//                persist_key, user_persist_key              , email               , verify_token01  ,register_date,verify_date, expire_date
		user.NewUserEmail(0 /*     */, savedUserRecords[0].PersistKey, "test01@example.com", "verify_token01", now /*    */, nil /*       */, nil),             // user不一致
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test02@example.com", "verify_token02", now /*    */, &now /*      */, nil),             // 更新対象
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test03@example.com", "verify_token03", now /*    */, nil /*       */, &now),            // verify_date null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test04@example.com", "verify_token04", now /*    */, &pastVerifyDate, &now),            // verify_date not null
		user.NewUserEmail(0 /*     */, savedUserRecords[1].PersistKey, "test05@example.com", "verify_token05", now /*    */, nil /*       */, &pastExpireDate), // expired
	}

	var expects = essence.ToVal(expectRecords)
	testUtility.AssertTable(t, orp, []string{"persist_key"}, expects, assertSameUserEmail)
}
