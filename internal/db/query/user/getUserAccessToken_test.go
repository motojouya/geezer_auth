package user_test

import (
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"github.com/motojouya/geezer_auth/internal/db/transfer/user"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetUser(t *testing.T) {
	var now = testUtility.GetNow()
	testUtility.Truncate(t, orp)

	var userRecords = []user.User{
		//           persist_key, identifier , email_idetifier     , name       , bot_flag  , register_date        , update_date
		user.NewUser(0 /*     */, "US-TASTAS", "test01@example.com", "tast name", false /**/, now.adddate(0, -1, 0), now.adddate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TESTES", "test01@example.com", "test name", false /**/, now.adddate(0, -1, 0), now.adddate(0, 0, -3)),
		user.NewUser(0 /*     */, "US-TOSTOS", "test01@example.com", "tost name", false /**/, now.adddate(0, -1, 0), now.adddate(0, 0, -3)),
	}
	var savedUserRecords = testUtility.Ready(t, orp, userRecords)

	var records = []user.UserAccessToken{
		//                      persist_key, user_persist_key              , access_token    , source_update_date   , register_date        , expire_date
		user.NewUserAccessToken(0 /*     */, savedUserRecords[0].PersistKey, "access_token05", now.adddate(0, 0, -3), now.AddDate(0, 0, -3), now.AddDate(0, 0, 7)), // user.identifierと一致しない
		user.NewUserAccessToken(0 /*     */, savedUserRecords[1].PersistKey, "access_token04", now.adddate(0, 0, -3), now.AddDate(0, 0, -4), now.AddDate(0, 0, 7)),
		user.NewUserAccessToken(0 /*     */, savedUserRecords[1].PersistKey, "access_token03", now.adddate(0, 0, -3), now.AddDate(0, 0, -5), now.AddDate(0, 0, 7)),
		user.NewUserAccessToken(0 /*     */, savedUserRecords[1].PersistKey, "access_token02", now.adddate(0, 0, -7), now.AddDate(0, 0, -6), now.AddDate(0, 0, 7)), // user.update_dateと一致しない
		user.NewUserAccessToken(0 /*     */, savedUserRecords[1].PersistKey, "access_token01", now.adddate(0, 0, -3), now.AddDate(0, 0, -7), now.AddDate(0, 0, -7)), // expire_dateが過去
	}
	var savedRecords = testUtility.Ready(t, orp, records)

	var result, err = orp.GetUser("US-TESTES")
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	// TODO expectの値を作る。2行目、3行目の2レコード検索できる想定

	assert.NotNil(t, result)
	assertSameUserAccessToken(t, savedRecords[1], *result)
}

func assertSameUserAccessToken(t *testing.T, expect user.UserAccessToken, actual user.UserAccessToken) {
	assert.Equal(t, expect.PersistKey, actual.PersistKey)
	assert.Equal(t, expect.Identifier, actual.Identifier)
	assert.Equal(t, expect.ExposeEmailId, actual.ExposeEmailId)
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.BotFlag, actual.BotFlag)
	assert.WithinDuration(t, expect.RegisteredDate, actual.RegisteredDate, time.Second)
	assert.WithinDuration(t, expect.UpdateDate, actual.UpdateDate, time.Second)
}
