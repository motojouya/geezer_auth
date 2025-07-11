package utility

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/go-gorp/gorp"
	"github.com/motojouya/geezer_auth/internal/core/essence"
)

var Dialect = goqu.Dialect("postgres")

/*
 * gorpのSelectOneは、レコードが見つからない場合にエラーになっちゃうがnilで返したいのでSelectSingleを作成
 */
func SelectSingle[R any](executer gorp.SqlExecutor, table string, keys map[string]string, query string, args ...interface{}) (*R, error) {
	var record []R
	var _, err = executer.Select(&record, query, args...)
	if err != nil {
		return nil, err
	}

	if len(record) == 0 {
		return nil, nil
	}

	if len(record) > 1 {
		return nil, essence.NewDuplicateError(table, keys, "Duplicate record found")
	}

	return &record[0], nil
}
