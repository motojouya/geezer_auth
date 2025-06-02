package authorization

import (
	"time"
	"errors"
)

/*
 * AuthorizationError
 */
type AuthorizationError struct {
	Role   string
	Action string
	error
}

func NewAuthorizationError(role, action, message string) *AuthorizationError {
	return &AuthorizationError{
		Role:    role,
		Action:  action,
		error:   errors.New(message),
	}
}

func (e AuthorizationError) Error() string {
	return e.error.Error() + " (role: " + e.Role + ", action: " + e.Action + ")"
}

func (e AuthorizationError) Unwrap() error {
	return e.error
}

func (e AuthorizationError) HttpStatus() uint {
	return 403
}

/*
 * TokenExpiredError
 *
 * JwtTokenが期限切れの場合に発生するエラー
 * Expireされたタイミングであっても、連続するsessionを継続しないと入力に1時間とかかかるやつのハンドリングが面倒なので、それを許容したい。
 * そのため、Expireでも形式的にエラーとしないため、jwtモジュールの外のAuthorizationに実装している。
 */
type TokenExpiredError struct {
	ExpiresAt time.Time
	error
}

func NewTokenExpiredError(expiresAt time.Time, message string) *TokenExpiredError {
	return &TokenExpiredError{
		ExpiresAt: expiresAt,
		error:     errors.New(message),
	}
}

func (e TokenExpiredError) Error() string {
	// TODO time.RFC3339 ?
	return e.error.Error() + " (expires at: " + e.ExpiresAt.Format(time.RFC3339) + ")"
}

func (e TokenExpiredError) Unwrap() error {
	return e.error
}

func (e TokenExpiredError) HttpStatus() uint {
	return 403
}
