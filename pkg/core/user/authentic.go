package user

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
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
	id UUID,
	user *User
) *Authentic {
	var expireDate = issueDate.Add(validityPeriodMinutes * time.Minute)

	return NewAuthentic(
		issuer,
		string(user.Identifier),
		audience,
		expiresAt,
		issuedAt,
		issuedAt,
		id,
		user
	)
}

func NewAuthentic(
	issuer    string,
	subject   string,
	audience  []string,
	expiresAt time.Time,
	notBefore time.Time,
	issuedAt  time.Time,
	id        UUID,
	user      *User
) *Authentic {
	return &Authentic{
		jwt.RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,      // iss
			Subject:   subject,     // sub
			Audience:  audience,    // aud
			ExpiresAt: expiresAt,   // exp
			NotBefore: notBefore,   // nbf
			IssuedAt:  issuedAt,    // iat
			ID:        id.String(), // jti
		},
		User: user,
	}
}
