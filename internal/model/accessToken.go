package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// TODO middlewareも作ってしまいたい。イメージを掴んで置く

type AccessToken struct {
	Token string
	ExpireDate time.Time
}

type Company struct {
	ExposeId string
	Name string
	Role string
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

// TODO expire dateを入れなくて大丈夫？
type GeezerClaims struct {
	jwt.RegisteredClaims
	UserEmail       *string   `json:"email"`
	UserName        string    `json:"name"`
	UpdateDate      time.Time `json:"update_at"`
	UserEmailId     string    `json:"email_id"` // TODO fqdnっぽくしないと、private keyとして文字列長さが短すぎる
	BotFlag         bool      `json:"bot_flag"` // TODO
	CompanyExposeId *string   `json:"company_expose_id"` // TODO
	CompanyName     *string   `json:"company_name"` // TODO
	CompanyRole     *string   `json:"company_role"` // TODO
}

// TODO GetId関数でuuid計算をDIできるようにしておく
type JwtSource struct {
	Issuer                string
	Audience              []string
	LatestSecretKeyId     string
	SecretMap             map[string]string
	ValidityPeriodMinutes uint
}

func NewJwtSource(issuer string, audience []string, latestSecret string, secretMap map[string]string, validityPeriodMinutes uint) JwtSource {
	return JwtSource{
		Issuer:                issuer,
		Audience:              audience,
		LatestSecret:          latestSecret,
		SecretMap:             secretMap,
		ValidityPeriodMinutes: validityPeriodMinutes,
	}
}

func CreateJwtSource() JwtSource {
	// TODO 環境変数とかかから作る
	return NewJwtSource(
		"issuer",
		[]string{"audience"},
		"latestSecret",
		map[string]string{"latestSecret":"latestSecret"},
		60,
	)
}

// expireの基準日がいるので、日付が必要。オプションで期間の調整ができてもいいかもしれない
func PublishAccessToken(jwtSource JwtSource, user User, tokens []AccessToken, currentTime time.Time) (AccessToken, error) {
	tokenString, err := jwtSource.GenerateAccessToken(user, currentTime)

	// TODO
	// 1. tokensをparseして、user.UpdateDateを取得する
	// 2. UpdateDateがuser.UpdateDateと違う場合は、GenerateAccessTokenを実行する
	// 3. 同じ場合でも、expireDateが過ぎている場合は、GenerateAccessTokenを実行する
	// 4. expireDateが過ぎていない場合でも、同一のUpdateDateのtokenでexpireDateがきれていないものが1つだけならば、GenerateAccessTokenを実行する
	// 5. そうでない場合は、tokensの中で最も新しいものをreturn

	return AccessToken{
		token: tokenString,
		expireDate: expireDate,
	}, nil
}


(jwtSource JwtSource) func GetExpireDate(issueDate time.Time) time.Time {
	return issueDate.Add(jwtSource.ValidityPeriodMinutes * time.Minute)
}

// TODO uuid発行してるので、参照透過じゃないねー。
// テストとしては、検証できるか否かになっちゃう？違うな、IDは検証できないだな。
// あるいは、IDを引数としていれるか
(jwtSource JwtSource) func GenerateAccessToken(user User, issueDate time.Time) (string, error) {
	var expireDate = jwtSource.GetExpireDate(issueDate)
	var claims = GeezerClaims{
		Issuer:    jwtSource.Issuer,               // iss
		Subject:   user.ExposeId,                  // sub
		Audience:  jwtSource.Audience,             // aud
		ExpiresAt: jwt.NewNumericDate(expireDate), // exp
		NotBefore: jwt.NewNumericDate(issueDate),  // nbf
		IssuedAt:  jwt.NewNumericDate(issueDate),  // iat
		ID:        uuid.NewUUID()                  // jti

		// TODO 他のproperty
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = jwtSource.LatestSecretKeyId

	var secret = jwtSource.SecretMap[jwtSource.LatestSecretKeyId]
	tokenString, err := token.SignedString(secret)
	if err != nil {
        	return "", err
	}

	return tokenString, nil
}

(jwtSource JwtSource) func GetUserFromAccessToken(jwtSecret JwtSecret, tokenString string) User {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		var secret = jwtSource.SecretMap[token.Header["kid"]]
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

 	token, err := jwt.ParseWithClaims(
 		tokenString,
 		&GeezerClaims{},
 		func(token *jwt.Token) (interface{}, error) {
			// var secret = jwtSource.SecretMap[token.Header["kid"]]
			// return secret, nil

 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
 				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
 			}
 			return secret, nil
 		},
 	)

	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(GeezerClaims); ok && token.Valid {
		return NewUser(claims.Issuer)
	} else {
		return nil
	}
}
