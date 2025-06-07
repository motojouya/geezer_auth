package text

import (
	pkg "github.com/motojouya/geezer_auth/pkg/core/text"
	"golang.org/x/crypto/bcrypt"
)

type Password string
type HashedPassword string

func NewPassword(password string) (Password, error) {
	if password == "" {
		return Password(""), pkg.NewLengthError("password", &password, 1, 255, "password should not be empty")
	}

	var length = len([]rune(password))
	if length < 1 || length > 255 {
		return Password(""), pkg.NewLengthError("password", &password, 1, 255, "password must be between 1 and 255 characters")
	}

	// TODO 正規表現あってる？
	re, err := regexp.Compile(`^[A-Za-z0-1]{1,255}$`)
	if err != nil {
		// 固定値なのでエラーにはならないはず
		panic(err)
	}

	var result = re.MatchString(text, -1)
	if !result {
		return Password(""), pkg.NewFormatError("password", "password", &password, "password must be a valid password")
	}

	return Password(password), nil
}

// DBから呼ばれる想定
func NewHashedPassword(password string) HashedPassword {
	return HashedPassword(password)
}

// 参照透過な関数ではないが、そもそもhashedPasswordはVerifyPassword関数を使わないと検証できず、VerifyPasswordと合わせると予測可能な動きとなるので、特にDIなどは必要ない
func HashPassword(password Password) (HashedPassword, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 9)
	if err != nil {
		return "", err
	}

	return HashedPassword(hashed), nil
}

func VerifyPassword(hashed HashedPassword, password Password) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(parameterPassword))
}
