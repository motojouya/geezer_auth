package user

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/google/uuid"
)

// TODO IsExpired(currentDate time.Time) bool 実装
type Authentic struct {
	jwt.RegisteredClaims
	User User
}

func CreateAuthentic(
	issuer string,
	audience []string,
	issuedAt time.Time,
	validityPeriodMinutes uint,
	id uuid.UUID,
	user *User,
) *Authentic {
	var expiresAt = issuedAt.Add(int(validityPeriodMinutes) * time.Minute)

	return NewAuthentic(
		issuer,
		string(user.Identifier),
		audience,
		expiresAt,
		issuedAt,
		issuedAt,
		id,
		*user,
	)
}

func NewAuthentic(
	issuer string,
	subject string,
	audience []string,
	expiresAt time.Time,
	notBefore time.Time,
	issuedAt time.Time,
	id uuid.UUID,
	user *User,
) *Authentic {
	return &Authentic{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,                        // iss
			Subject:   subject,                       // sub
			Audience:  audience,                      // aud
			ExpiresAt: jwt.NewNumericDate(expiresAt), // exp
			NotBefore: jwt.NewNumericDate(notBefore), // nbf
			IssuedAt:  jwt.NewNumericDate(issuedAt),  // iat
			ID:        id.String(),                   // jti
		},
		User: *user,
	}
}
