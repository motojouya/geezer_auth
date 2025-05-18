package accessToken

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// TODO IsExpired(currentDate time.Time) bool 実装
type GeezerToken struct {
	Issuer    string
	Subject   string
	Audience  []string
	ExpiresAt time.Time
	NotBefore time.Time
	IssuedAt  time.Time
	ID        string
	User      User
}

func NewGeezerToken(issuer string, subject string, audience []string, expiresAt time.Time, notBefore time.Time, issuedAt time.Time, id string, user User) *GeezerToken {
	return &GeezerToken{
		Issuer:    issuer,
		Subject:   subject,
		Audience:  audience,
		ExpiresAt: expiresAt,
		NotBefore: notBefore,
		IssuedAt:  issuedAt,
		ID:        id,
		User:      user,
	}
}

// FIXME claimsのprivate keyが`github.com/motojouya/geezer_auth/`をprefixとしているが、本来は稼働するサーバのfqdnをprefixとして持つべき。
type GeezerClaims struct {
	jwt.RegisteredClaims
	UserEmail       *string   `json:"email"`
	UserName        string    `json:"name"`
	UpdateDate      time.Time `json:"update_at"`
	UserEmailId     string    `json:"github.com/motojouya/geezer_auth/email_id"`
	BotFlag         bool      `json:"github.com/motojouya/geezer_auth/bot_flag"`
	CompanyExposeId *string   `json:"github.com/motojouya/geezer_auth/company_expose_id"`
	CompanyName     *string   `json:"github.com/motojouya/geezer_auth/company_name"`
	CompanyRole     *string   `json:"github.com/motojouya/geezer_auth/company_role"`
}

func CreateClaims(user User, issuer string, audience []string, expiresAt time.Time, issuedAt time.Time, id string) *GeezerClaims {
	return GeezerClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,                         // iss
			Subject:   user.ExposeId,                  // sub
			Audience:  audience,                       // aud
			ExpiresAt: jwt.NewNumericDate(expireDate), // exp
			NotBefore: jwt.NewNumericDate(issueDate),  // nbf
			IssuedAt:  jwt.NewNumericDate(issueDate),  // iat
			ID:        id                              // jti
		},
		UserEmail:       user.Email,
		UserName:        user.Name,
		UpdateDate:      jwt.NewNumericDate(user.UpdateDate),
		UserEmailId:     user.EmailId,
		BotFlag:         user.BotFlag,
		CompanyExposeId: user.Company.ExposeId,
		CompanyName:     user.Company.Name,
		CompanyRole:     user.Company.Role,
	}
}
