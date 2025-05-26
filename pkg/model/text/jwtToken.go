package model

import (
	"time"
)

type JwtToken string

// DBから来た値のみを想定。Generateする際はpkg/accessTokenに任せる
func NewJwtToken(token string) JwtToken {
	return JwtToken(token)
}

// TODO これそもそもuser_updated_atが、DBに入ってたら、それと照合してDBだけで条件検索できる。ので、単純にgenerateだけで良い気がしてきた。
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
