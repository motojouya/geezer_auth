package command_test

import (
	"github.com/motojouya/geezer_auth/internal/db/query/command"
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/shelter/essence"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddPassword(t *testing.T) {
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

	var pastExpireDate = now.AddDate(0, -1, 0)
	var records = []*user.UserPassword{
		//                   persist_key,user_persist_key, password    ,register_date,expire_date
		user.NewUserPassword(0 /*     */, savedUserRecords[0].PersistKey, "password01", now /*    */, nil),             // user 不一致
		user.NewUserPassword(0 /*     */, savedUserRecords[1].PersistKey, "password02", now /*    */, nil),             // 現行
		user.NewUserPassword(0 /*     */, savedUserRecords[1].PersistKey, "password03", now /*    */, &pastExpireDate), // expire
	}
	testUtility.ReadyPointer(t, orp, records)

	var newRecord = user.NewUserPassword(0, savedUserRecords[1].PersistKey, "password04", now.AddDate(0, 0, 3), nil)
	var savedRecord, err = command.AddPassword(orp, newRecord, now)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	assert.NotNil(t, savedRecord)
	assertSameUserPassword(t, *newRecord, *savedRecord)

	var expectRecords = []*user.UserPassword{
		//                   persist_key, user_persist_key              , password    , register_date       , expire_date
		user.NewUserPassword(1 /*     */, savedUserRecords[0].PersistKey, "password01", now /*            */, nil),             // user 不一致
		user.NewUserPassword(2 /*     */, savedUserRecords[1].PersistKey, "password02", now /*            */, &now),            // 旧
		user.NewUserPassword(3 /*     */, savedUserRecords[1].PersistKey, "password03", now /*            */, &pastExpireDate), // expire
		user.NewUserPassword(4 /*     */, savedUserRecords[1].PersistKey, "password04", now.AddDate(0, 0, 3), nil),             // 新規
	}

	var expects = essence.ToVal(expectRecords)
	testUtility.AssertTable(t, orp, []string{"persist_key"}, expects, assertSameUserPassword)
}
