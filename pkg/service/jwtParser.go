package service

import (
	"github.com/motojouya/geezer_auth/pkg/core/jwt"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"github.com/motojouya/geezer_auth/pkg/utility"
	"os"
)

type JwtParserLoader interface {
	LoadJwtHandler(local io.Local) (JwtHandler, error)
}

type jwtParserLoaderImpl interface{}

type JwtParser interface {
	Parse(tokenString string) (*user.Authentic, error)
}

// FIXME 本来はinternal/io/localのGetEnvを使って環境変数を取得すべきだが、ioパッケージを共通部品として使うのは不自然なので、ここではosパッケージを直接使う形
// また、osパッケージ自体がローカルマシンに依存したものではあるので、引数で受け取ってコントローラブルにすべき
func (impl jwtParserLoaderImpl) LoadJwtParser() (JwtParser, error) {
	if issuer, issuerExist := os.LookupEnv("JWT_ISSUER"); !issuerExist {
		return nil, utility.NewSystemConfigError("JWT_ISSUER", "JWT_ISSUER is not set on env")
	}
	if myself, myselfExist := os.LookupEnv("JWT_MYSELF"); !myselfExist {
		return nil, utility.NewSystemConfigError("JWT_ISSUER", "JWT_ISSUER is not set on env")
	}
	if latestKeyId, latestKeyIdExist := os.LookupEnv("JWT_LATEST_KEY_ID"); !latestKeyIdExist {
		return nil, utility.NewSystemConfigError("JWT_LATEST_KEY_ID", "JWT_LATEST_KEY_ID is not set on env")
	}
	if latestSecret, latestSecretExist := os.LookupEnv("JWT_LATEST_SECRET"); !latestSecretExist {
		return nil, utility.NewSystemConfigError("JWT_LATEST_SECRET", "JWT_LATEST_SECRET is not set on env")
	}
	if oldKeyId, oldKeyIdExist := os.LookupEnv("JWT_OLD_KEY_ID"); !oldKeyIdExist {
		return nil, utility.NewSystemConfigError("JWT_OLD_KEY_ID", "JWT_OLD_KEY_ID is not set on env")
	}
	if oldSecret, oldSecretExist := os.LookupEnv("JWT_OLD_SECRET"); !oldSecretExist {
		return nil, utility.NewSystemConfigError("JWT_OLD_SECRET", "JWT_OLD_SECRET is not set on env")
	}

	return jwt.NewJwtParsing(
		issuer,
		myself,
		latestKeyId,
		latestSecret,
		oldKeyId,
		oldSecret,
	), nil
}
