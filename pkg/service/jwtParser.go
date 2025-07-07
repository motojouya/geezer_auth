package service

import (
	"github.com/motojouya/geezer_auth/pkg/core/essence"
	"github.com/motojouya/geezer_auth/pkg/core/jwt"
	"github.com/motojouya/geezer_auth/pkg/core/user"
	"os"
)

type JwtParserLoader interface {
	LoadJwtParser() (JwtParser, error)
}

type jwtParserLoaderImpl struct{}

type JwtParser interface {
	Parse(tokenString string) (*user.Authentic, error)
}

var jwtParser *jwt.JwtParsing

// FIXME 本来はinternal/io/localのGetEnvを使って環境変数を取得すべきだが、ioパッケージを共通部品として使うのは不自然なので、ここではosパッケージを直接使う形
// また、osパッケージ自体がローカルマシンに依存したものではあるので、引数で受け取ってコントローラブルにすべき
func (impl jwtParserLoaderImpl) LoadJwtParser() (JwtParser, error) {
	if jwtParser != nil {
		return jwtParser, nil
	}

	var issuer, issuerExist = os.LookupEnv("JWT_ISSUER")
	if !issuerExist {
		return nil, essence.NewSystemConfigError("JWT_ISSUER", "JWT_ISSUER is not set on env")
	}

	var myself, myselfExist = os.LookupEnv("JWT_MYSELF")
	if !myselfExist {
		return nil, essence.NewSystemConfigError("JWT_ISSUER", "JWT_ISSUER is not set on env")
	}

	var latestKeyId, latestKeyIdExist = os.LookupEnv("JWT_LATEST_KEY_ID")
	if !latestKeyIdExist {
		return nil, essence.NewSystemConfigError("JWT_LATEST_KEY_ID", "JWT_LATEST_KEY_ID is not set on env")
	}

	var latestSecret, latestSecretExist = os.LookupEnv("JWT_LATEST_SECRET")
	if !latestSecretExist {
		return nil, essence.NewSystemConfigError("JWT_LATEST_SECRET", "JWT_LATEST_SECRET is not set on env")
	}

	var oldKeyId, oldKeyIdExist = os.LookupEnv("JWT_OLD_KEY_ID")
	if !oldKeyIdExist {
		return nil, essence.NewSystemConfigError("JWT_OLD_KEY_ID", "JWT_OLD_KEY_ID is not set on env")
	}

	var oldSecret, oldSecretExist = os.LookupEnv("JWT_OLD_SECRET")
	if !oldSecretExist {
		return nil, essence.NewSystemConfigError("JWT_OLD_SECRET", "JWT_OLD_SECRET is not set on env")
	}

	var jwtParsing = jwt.NewJwtParsing(
		issuer,
		myself,
		latestKeyId,
		latestSecret,
		oldKeyId,
		oldSecret,
	)

	jwtParser = &jwtParsing

	return jwtParser, nil
}
