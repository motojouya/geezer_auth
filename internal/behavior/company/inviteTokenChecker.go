package user

import (
	userQuery "github.com/motojouya/geezer_auth/internal/db/query/user"
	entryAuth "github.com/motojouya/geezer_auth/internal/entry/transfer/auth"
	localPkg "github.com/motojouya/geezer_auth/internal/local"
	shelterUser "github.com/motojouya/geezer_auth/internal/shelter/user"
)

type RefreshTokenCheckerDB interface {
	userQuery.GetUserRefreshTokenQuery
}

type RefreshTokenChecker interface {
	Execute(entry entryAuth.RefreshTokenGetter) (*shelterUser.UserAuthentic, error)
}

type RefreshTokenCheck struct {
	local localPkg.Localer
	db    RefreshTokenCheckerDB
}

func NewRefreshTokenCheck(local localPkg.Localer, database RefreshTokenCheckerDB) *RefreshTokenCheck {
	return &RefreshTokenCheck{
		db:    database,
		local: local,
	}
}

func (checker RefreshTokenCheck) Execute(entry entryAuth.RefreshTokenGetter) (*shelterUser.UserAuthentic, error) {
	now := checker.local.GetNow()

	refreshToken, err := entry.GetRefreshToken()
	if err != nil {
		return nil, err
	}

	// FIXME refresh tokenは単体で認証もできるので漏れるとだいぶまずい。passwordと同様に不可逆暗号として保存し、検索時に暗号化して当てるべき。
	dbUserAuthentic, err := checker.db.GetUserRefreshToken(string(refreshToken), now)
	if err != nil {
		return nil, err
	}

	return dbUserAuthentic.ToShelterUserAuthentic()
}
