package testUtility

import (
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"reflect"
	"testing"
	"github.com/stretchr/testify/assert"
)

func Truncate(t *testing.T, orp db.ORP) {
	var impl, ok = orp.(*db.ORPImpl)
	if !ok {
		t.Fatalf("Expected db.ORPImpl, got %T", orp)
	}
	var err = impl.DbMap.TruncateTables()
	if err != nil {
		t.Fatalf("Could not truncate tables: %s", err)
	}
}

func Ready[T any](t *testing.T, orp db.ORP, records ...T) {
	var r []interface{}
	for _, record := range records {
		r = append(r, record)
	}
	var err = orp.Insert(r...)
	if err != nil {
		t.Fatalf("Could not insert records: %s", err)
	}
}

func AssertTable[T any](t *testing.T, orp db.ORP, expects ...T) {
	var impl, ok = orp.(*db.ORPImpl)
	if !ok {
		t.Fatalf("Expected db.ORPImpl, got %T", orp)
	}

	var zero T
	var table, tableErr = impl.DbMap.TableFor(reflect.TypeOf(zero).Elem(), true)
	if tableErr != nil {
		t.Fatalf("Could not get table for %T: %s", zero, tableErr)
	}

	// FIXME gorpがtableまでは判明させてくれるが、主キーは取得できない。
	// `assert.ElementsMatch`が順序関係なく照合してくれるのでorder by句は実質不要だが、なんか気持ち悪いところ
	// var orders []exp.OrderedExpression
	// for _, column := range table.Columns {
	// 	if column.isPK {
	// 		orders = append(orders, goqu.C(column.ColumnName).Asc())
	// 	}
	// }
	// var sql, args, sqlErr = utility.Dialect.From(table.TableName).Order(orders...).ToSQL()
	var sql, args, sqlErr = utility.Dialect.From(table.TableName).ToSQL()
	if sqlErr != nil {
		t.Fatalf("Could not create SQL for %T: %s", zero, sqlErr)
	}

	var actuals []T
	var _, execErr = impl.Select(&actuals, sql, args...)
	if execErr != nil {
		t.Fatalf("Could not execute SQL for %T: %s", zero, execErr)
	}

	assert.ElementsMatchf(t, expects, actuals, "Expected records do not match actual records in table %s", table.TableName)
}
