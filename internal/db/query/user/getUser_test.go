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

	var records = []user.User{
		//           persist_key,identifier ,email_idetifier     ,name       ,bot_flag ,,register_date,update_date
		user.NewUser(0/*      */,"US-TASTAS","test01@example.com","tast name",false/**/,now/*       */,now),
		user.NewUser(0/*      */,"US-TESTES","test01@example.com","test name",false/**/,now/*       */,now),
		user.NewUser(0/*      */,"US-TOSTOS","test01@example.com","tost name",false/**/,now/*       */,now),
	}

	testUtility.Truncate(t, orp)
	var savedRecords = testUtility.Ready(t, orp, records)

	var result, err = orp.GetUser("US-TESTES")
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}

	assert.NotNil(t, result)
	assertSameUser(t, savedRecords[1], *result)
}

func assertSameUser(t *testing.T, expect user.User, actual user.User) {
	assert.Equal(t, expect.PersistKey, actual.PersistKey)
	assert.Equal(t, expect.Identifier, actual.Identifier)
	assert.Equal(t, expect.ExposeEmailId, actual.ExposeEmailId)
	assert.Equal(t, expect.Name, actual.Name)
	assert.Equal(t, expect.BotFlag, actual.BotFlag)
	assert.WithinDuration(t, expect.RegisteredDate, actual.RegisteredDate, time.Second)
	assert.WithinDuration(t, expect.UpdateDate, actual.UpdateDate, time.Second)
}
