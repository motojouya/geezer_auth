package model

import (
	"time"
	"fmt"
)

// Nameは長さ10000文字までのstringを表す型
// DBに保存する際にレコードと同一ページに収まるぐらいのサイズを想定している。 TODO 10000文字が妥当か要検討
// これより長い文字列が必要な場合は、別途LongText型を作成すること。LongText型はDBに直接配置するのではなく、Object Storageに保存することを想定すべき
type Text *string

func NewText(text *string) (*Text, error) {
	var length = len(*[]rune(name))
	if length > 10000 {
		return Name(""), fmt.Errorf("name must be less then 10000 characters")
	}

	return &Text(text), nil
}
