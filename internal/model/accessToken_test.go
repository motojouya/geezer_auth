package model_test

import (
	"github.com/motojouya/geezer_auth/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewAccessToken(t *testing.T) {
	var token = "test_token"
	var expireDate = time.Now()

	var accessToken = model.NewAccessToken(token, expireDate)

	assert.Equal(t, token, accessToken.Token)
	assert.Equal(t, expireDate, accessToken.ExpireDate)

	t.Logf("accessToken: %+v", accessToken)
	t.Logf("accessToken.Token: %d", role.Token)
	t.Logf("accessToken.ExpireDate: %s", role.ExpireDate)
}
