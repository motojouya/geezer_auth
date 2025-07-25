package text

import (
	pkg "github.com/motojouya/geezer_auth/pkg/core/text"
	"strings"
)

/*
 * Textは長さ一万文字までのstringを表す型
 * MySQLのTEXT型に相当するが、MySQLのTEXT型は65535バイトまで格納でき、utf8の4byte文字を格納するとだいたい一万文字強まで格納できるため。
 * 本プロジェクトはMySQLではなくPostgreSQLを想定しているが、ちょうど良いので採用した。千文字では少し足りないし、二万文字では多すぎる。
 * これより長い文字列が必要な場合は、別途LongText型を作成すること。LongText型はDBに直接配置するのではなく、Object Storageに保存することを想定すべき
 */
type Text string

func NewText(text string) (Text, error) {
	var trimmed = strings.TrimSpace(text)

	var length = len([]rune(trimmed))
	if length > 10000 {
		return Text(""), pkg.NewLengthError("text", text, 0, 10000, "text must be between 0 and 10000 characters")
	}

	return Text(trimmed), nil
}
