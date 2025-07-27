package command_test

import (
	"github.com/motojouya/geezer_auth/internal/core/essence"
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/motojouya/geezer_auth/internal/db/query/command"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExpirePassword(t *testing.T) {
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
		//                   persist_key, user_persist_key              , password    ,register_date,expire_date
		user.NewUserPassword(0 /*     */, savedUserRecords[0].PersistKey, "password01", now /*    */, nil),             // x user 不一致
		user.NewUserPassword(0 /*     */, savedUserRecords[1].PersistKey, "password02", now /*    */, nil),             // o 対象
		user.NewUserPassword(0 /*     */, savedUserRecords[1].PersistKey, "password03", now /*    */, &pastExpireDate), // x expire
	}
	testUtility.ReadyPointer(t, orp, records)

	var err = command.ExpirePassword(orp, savedUserRecords[1].PersistKey, now)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	var expectRecords = []*user.UserPassword{
		//                   persist_key, user_persist_key              , password    ,register_date,expire_date
		user.NewUserPassword(1 /*     */, savedUserRecords[0].PersistKey, "password01", now /*    */, nil),             // x user 不一致
		user.NewUserPassword(2 /*     */, savedUserRecords[1].PersistKey, "password02", now /*    */, &now),            // o 対象
		user.NewUserPassword(3 /*     */, savedUserRecords[1].PersistKey, "password03", now /*    */, &pastExpireDate), // x expire
	}

	var expects = essence.ToVal(expectRecords)
	testUtility.AssertTable(t, orp, []string{"persist_key"}, expects, assertSameUserPassword)
}

func assertSameUserPassword(t *testing.T, expect user.UserPassword, actual user.UserPassword) {
	assert.Equal(t, expect.UserPersistKey, actual.UserPersistKey)
	assert.Equal(t, expect.Password, actual.Password)
	assert.WithinDuration(t, expect.RegisteredDate, actual.RegisteredDate, time.Second)
	if expect.ExpireDate == nil {
		assert.Nil(t, actual.ExpireDate)
	} else {
		assert.WithinDuration(t, *expect.ExpireDate, *actual.ExpireDate, time.Second)
	}
}
