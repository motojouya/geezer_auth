package command_test

import (
	"github.com/motojouya/geezer_auth/internal/db/query/command"
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"testing"
)

func TestAddRefreshToken(t *testing.T) {
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

	var records = []user.UserRefreshToken{
		//                       persist_key, user_persist_key              , refresh_token    , register_date        , expire_date
		user.NewUserRefreshToken(0 /*     */, savedUserRecords[0].PersistKey, "refresh_token01", now.AddDate(0, -1, 0), now.AddDate(0, 0, 7)),  // x user不一致
		user.NewUserRefreshToken(0 /*     */, savedUserRecords[1].PersistKey, "refresh_token02", now.AddDate(0, -1, 0), now.AddDate(0, 0, 7)),  // o 対象
		user.NewUserRefreshToken(0 /*     */, savedUserRecords[1].PersistKey, "refresh_token03", now.AddDate(0, -1, 0), now.AddDate(0, 0, -7)), // x expire_date過去
	}
	testUtility.Ready(t, orp, records)

	var newRecord = user.NewUserRefreshToken(0, savedUserRecords[1].PersistKey, "refresh_token04", now, now.AddDate(0, 0, 9))
	var resultRecord, err = command.AddRefreshToken(orp, newRecord, now)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	assertSameUserRefreshToken(t, newRecord, resultRecord)

	var expectRecords = []user.UserRefreshToken{
		//                       persist_key, user_persist_key              , refresh_token    , register_date        , expire_date
		user.NewUserRefreshToken(0 /*     */, savedUserRecords[0].PersistKey, "refresh_token01", now.AddDate(0, -1, 0), now.AddDate(0, 0, 7)),  // user不一致
		user.NewUserRefreshToken(0 /*     */, savedUserRecords[1].PersistKey, "refresh_token02", now.AddDate(0, -1, 0), now),                   // 旧
		user.NewUserRefreshToken(0 /*     */, savedUserRecords[1].PersistKey, "refresh_token03", now.AddDate(0, -1, 0), now.AddDate(0, 0, -7)), // expire_date過去
		user.NewUserRefreshToken(0 /*     */, savedUserRecords[1].PersistKey, "refresh_token04", now /*             */, now.AddDate(0, 0, 9)),  // 新規追加
	}

	testUtility.AssertTable(t, orp, []string{"persist_key"}, expectRecords, assertSameUserRefreshToken)
}
