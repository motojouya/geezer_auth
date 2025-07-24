package testUtility

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"reflect"
	"testing"
	"time"
)

func GetNow() time.Time {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	return time.Now().In(jst)
}

// truncateの順序を決めているので、順番が重要。依存される側があとに来るようにする。
var tables = []string{
	"user_access_token",
	//"user_refresh_token",
	//"user_password",
	"user_company_role",
	"user_email",
	"company_invite",
	"users",
	"company",
	"role_permission",
	"role",
}

func Truncate(t *testing.T, orp db.ORP) {
	for _, table := range tables {
		var _, err = orp.Exec("TRUNCATE TABLE " + table + " CASCADE")
		if err != nil {
			t.Fatalf("Could not truncate table %s: %s", table, err)
		}
	}
}

// truncateはforeign key制約のためにcascadeをつける必要があるので、独自実装で行っている。
func oldTruncate(t *testing.T, orp db.ORP) {
	var impl, ok = orp.(*db.ORPImpl)
	if !ok {
		t.Fatalf("Expected db.ORPImpl, got %T", orp)
	}
	var err = impl.DbMap.TruncateTables()
	if err != nil {
		t.Fatalf("Could not truncate tables: %s", err)
	}
}

func Ready[T any](t *testing.T, orp db.ORP, records []T) []T {
	var rec []interface{}
	for _, record := range records {
		rec = append(rec, &record)
	}
	var err = orp.Insert(rec...)
	if err != nil {
		t.Fatalf("Could not insert records: %s", err)
	}

	var ret = make([]T, len(records))
	for i, r := range rec {
		var result, ok = r.(*T)
		if !ok {
			t.Fatalf("Expected type %T, got %T", ret[i], r)
		}
		ret[i] = *result
	}
	return ret
}

func ReadyPointer[T any](t *testing.T, orp db.ORP, records []*T) []*T {
	var rec []interface{}
	for _, record := range records {
		rec = append(rec, record)
	}
	var err = orp.Insert(rec...)
	if err != nil {
		t.Fatalf("Could not insert records: %s", err)
	}

	var ret = make([]*T, len(records))
	for i, r := range rec {
		var result, ok = r.(*T)
		if !ok {
			t.Fatalf("Expected type %T, got %T", ret[i], r)
		}
		ret[i] = result
	}
	return ret
}

func AssertRecords[T any](t *testing.T, expects []T, actuals []T, assertSame func(*testing.T, T, T)) {
	if len(expects) != len(actuals) {
		t.Fatalf("Expected %d records, got %d", len(expects), len(actuals))
	}

	for i := len(expects) - 1; i >= 0; i-- {
		assertSame(t, expects[i], actuals[i])
	}
}

func AssertTable[T any](t *testing.T, orp db.ORP, orders []string, expects []T, assertSame func(*testing.T, T, T)) {
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
	// var orderBys []exp.OrderedExpression
	// for _, column := range table.Columns {
	// 	if column.isPK {
	// 		orderBys = append(orderBys, goqu.C(column.ColumnName).Asc())
	// 	}
	// }
	// var sql, args, sqlErr = utility.Dialect.From(table.TableName).Order(orderBys...).ToSQL()
	var orderBys []exp.OrderedExpression
	for _, order := range orders {
		orderBys = append(orderBys, goqu.C(order).Asc())
	}

	var sql, args, sqlErr = utility.Dialect.From(table.TableName).Order(orderBys...).ToSQL()
	if sqlErr != nil {
		t.Fatalf("Could not create SQL for %T: %s", zero, sqlErr)
	}

	var actuals []T
	var _, execErr = impl.Select(&actuals, sql, args...)
	if execErr != nil {
		t.Fatalf("Could not execute SQL for %T: %s", zero, execErr)
	}

	AssertRecords(t, expects, actuals, assertSame)
}
