package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"os"
	"strconv"
)

// TODO 本packageは、externalにするので、どこからも依存しない感じにしておく。
// 変換としては、model.Userから、このuserへの変換が必要だが、それはmodelに生やす感じで

type Company struct {
	ExposeId string
	Name string
	Role string
}

func NewCompany(exposeId string, name string, role string) *Company {
	return &Company{
		ExposeId: exposeId,
		Name: name,
		Role: role,
	}
}

type User struct {
	ExposeId string
	EmailId string
	Email *string
	Name string
	BotFlag bool
	Company *Company
	UpdateDate time.Time
}

func NewUser(exposeId string, emailId string, email *string, name string, botFlag bool, company *Company, updateDate time.Time) *User {
	return &User{
		ExposeId: exposeId,
		EmailId: emailId,
		Email: email,
		Name: name,
		BotFlag: botFlag,
		Company: company,
		UpdateDate: updateDate,
	}
}

type GeezerToken struct {
	Issuer    string
	Subject   string
	Audience  []string
	ExpiresAt time.Time
	NotBefore time.Time
	IssuedAt  time.Time
	ID        string
	User      User
	Valid     bool
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

func CreateGeezerToken(claims GeezerClaims) (*GeezerToken, error) {
	var company *Company = nil
	if claims.CompanyExposeId != nil && claims.CompanyName != nil && claims.CompanyRole != nil {
		company = NewCompany(*claims.CompanyExposeId, *claims.CompanyName, *claims.CompanyRole)
	} else {
		if claims.CompanyExposeId != nil {
			return nil, fmt.Error("CompanyExposeId is not nil")
		}
		if claims.CompanyName != nil {
			return nil, fmt.Error("CompanyName is not nil")
		}
		if claims.CompanyRole != nil {
			return nil, fmt.Error("CompanyRole is not nil")
		}
		// company = nil
	}

	var user = NewUser(
		claims.Subject,
		claims.UserEmailId,
		claims.UserEmail,
		claims.UserName,
		claims.BotFlag,
		company,
		claims.UpdateDate, // TODO 値としてtime.Time型になってる？
	)
}

type JwtParser struct {
	Issuer            string
	Myself            string // Audienceと付き合わせるための自分自身の情報
	LatestSecretKeyId string
	SecretMap         map[string]string
}

type JwtHandler struct {
	Audience              []string
	ValidityPeriodMinutes uint
	GetId                 func () (string, error)
	JwtParser
}

func NewJwtHandler(audience []string, jwtParser JwtParser, validityPeriodMinutes uint, GetId (func() (string, error))) *JwtHandler {
	return &JwtHandler{
		Issuer:                issuer,
		Audience:              audience,
		ValidityPeriodMinutes: validityPeriodMinutes,
		GetId:                 GetId,
		JwtParser:             jwtParser
	}
}

func NewJwtParser(issuer string, myself string, latestSecret string, secretMap map[string]string) *JwtParser {
	return &JwtParser{
		Issuer:       issuer,
		Myself:       myself,
		LatestSecret: latestSecret,
		SecretMap:    secretMap,
	}
}


func CreateJwtAudienceParser() (*JwtParser, error) {
	if issuer, issuerExist := os.LookupEnv("JWT_ISSUER"); !issuerExist {
		return nil, fmt.Error("JWT_ISSUER is not set on env")
	}
	if secretKeyId, secretKeyIdExist := os.LookupEnv("JWT_SECRET_KEY_ID"); !secretKeyIdExist {
		return nil, fmt.Error("JWT_SECRET_KEY_ID is not set on env")
	}
	if secret, secretExist := os.LookupEnv("JWT_SECRET"); !secretExist {
		return nil, fmt.Error("JWT_SECRET is not set on env")
	}
	if myself, myselfExist := os.LookupEnv("JWT_MYSELF"); !myselfExist {
		return nil, fmt.Error("JWT_ISSUER is not set on env")
	}

	return NewJwtParser(
		issuer,
		myself,
		secretKeyId,
		map[string]string{secretKeyId:secret},
	), nil
}

func CreateJwtIssuerParser() (*JwtParser, error) {
	if issuer, issuerExist := os.LookupEnv("JWT_ISSUER"); !issuerExist {
		return nil, fmt.Error("JWT_ISSUER is not set on env")
	}
	if secretKeyId, secretKeyIdExist := os.LookupEnv("JWT_SECRET_KEY_ID"); !secretKeyIdExist {
		return nil, fmt.Error("JWT_SECRET_KEY_ID is not set on env")
	}
	if secret, secretExist := os.LookupEnv("JWT_SECRET"); !secretExist {
		return nil, fmt.Error("JWT_SECRET is not set on env")
	}

	return NewJwtParser(
		issuer,
		issuer,
		secretKeyId,
		map[string]string{secretKeyId:secret},
	), nil
}

func CreateJwtHandler() (*JwtHandler, error) {
	if audience, audienceExist := os.LookupEnv("JWT_AUDIENCE"); !audienceExist {
		return nil, fmt.Error("JWT_AUDIENCE is not set on env")
	}
	if validityPeriodMinutesStr, validityPeriodMinutesExist := os.LookupEnv("JWT_VALIDITY_PERIOD_MINUTES"); !validityPeriodMinutesExist {
		return nil, fmt.Error("JWT_VALIDITY_PERIOD_MINUTES is not set on env")
	}
	var validityPeriodMinutes, err = strconv.Atoi(validityPeriodMinutesStr)
	if err != nil {
		return nil, err
	}

	var GetId = func () (string, error) {
		token, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}

		return token.String(), nil
	}

	var jwtParser, err = CreateJwtIssuerParser()
	if err != nil {
		return nil, err
	}
	if jwtParser == nil {
		return nil, fmt.Error("jwtParser is nil")
	}

	return NewJwtHandler(
		issuer,
		[]string{audience},
		*jwtParser,
		validityPeriodMinutes,
		GetId,
	), nil
}

func CreateClaims(jwtHandler JwtHandler, user User, issueDate time.Time) GeezerClaims {

	var id, err = jwtSource.GetId()
	if err != nil {
		return "", err
	}

	var expireDate = issueDate.Add(jwtSource.ValidityPeriodMinutes * time.Minute)

	return GeezerClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtSource.Issuer,               // iss
			Subject:   user.ExposeId,                  // sub
			Audience:  jwtSource.Audience,             // aud
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

(jwtHandler JwtHandler) func GenerateAccessToken(user User, issueDate time.Time) (string, error) {
	var claims = CreateClaims(jwtSource, user, issueDate)

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = jwtSource.LatestSecretKeyId

	var secret = jwtSource.SecretMap[jwtSource.LatestSecretKeyId]
	tokenString, err := token.SignedString(secret)
	if err != nil {
        	return "", err
	}

	return tokenString, nil
}

// TODO ExpiresAtの時間を取得する issuerとかの情報とかも見たほうがいい
(jwtParser JwtParser) func GetUserFromAccessToken(tokenString string) (*GeezerToken, error) {
 	token, err := jwt.ParseWithClaims(
 		tokenString,
 		&GeezerClaims{},
 		func(token *jwt.Token) (interface{}, error) {
			// jwt.SigningMethodHMAC?
 			if _, ok := token.Method.(*jwt.SigningMethodHS256); !ok {
 				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
 			}

			var secret, exist = jwtParser.SecretMap[token.Header["kid"]]
			if !exist {
				return nil, fmt.Errorf("Secret not found for key: %v", token.Header["kid"])
			}
			return secret, nil
 		},
 	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(GeezerClaims); !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return CreateGeezerToken(claims)
}

// TODO middlewareも作ってしまいたい。イメージを掴んで置く
// 
// type AccessToken struct {
// 	Token string
// 	ExpireDate time.Time
// }
// 
// これそもそもuser_updated_atが、DBに入ってたら、それと照合してDBだけで条件検索できる。ので、単純にgenerateだけで良い気がしてきた。
// // expireの基準日がいるので、日付が必要。オプションで期間の調整ができてもいいかもしれない
// func PublishAccessToken(jwtHandler JwtHandler, user User, tokens []AccessToken, currentTime time.Time) (AccessToken, error) {
// 	tokenString, err := jwtSource.GenerateAccessToken(user, currentTime)
// 
// 	// TODO
// 	// 1. tokensをparseして、user.UpdateDateを取得する
// 	// 2. UpdateDateがuser.UpdateDateと違う場合は、GenerateAccessTokenを実行する。一致を見るのに、誤差は許容したほうがいいかもしれない
// 	// 3. 同じ場合でも、expireDateが過ぎている場合は、GenerateAccessTokenを実行する
// 	// 4. expireDateが過ぎていない場合でも、同一のUpdateDateのtokenでexpireDateがきれていないものが1つだけならば、GenerateAccessTokenを実行する
// 	// 5. そうでない場合は、tokensの中で最も新しいものをreturn
// 
// 	return AccessToken{
// 		token: tokenString,
// 		expireDate: expireDate,
// 	}, nil
// }
// (jwtHandler JwtHandler) func GetExpireDate(issueDate time.Time) time.Time {
// 	return issueDate.Add(jwtSource.ValidityPeriodMinutes * time.Minute)
// }
