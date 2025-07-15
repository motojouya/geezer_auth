package role_test

import (
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/motojouya/geezer_auth/internal/db/testUtility"
	"testing"
)

var orp db.ORP

func TestMain(m *testing.M) {
	testUtility.ExecuteDatabaseTest(func(orpArg db.ORP) int {
		orp = orpArg
		return m.Run()
	})
	orp = nil // il?
}

// TODO transaction管理どうしようか
// func TestGreatestUsecase(t *testing.T) {
// 	tests := map[string]struct {
// 		input   GreatestUsecaseInput
// 		want    GreatestUsecaseOutput
// 	}{
// 		"when greatest usecase": {input: GreatestUsecaseInput{Something: "師走"}, want: GreatestUsecaseInput{Something: "師走"}},
// 	}
//
// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			var output GreatestUsecaseOuntput
// 			greatestUsecase := NewGreatestUsecase(repository.NewGreatestRepository())
// 			ctx := context.Background()
//
//            // NOTE: テストケースごとにトランザクションを開始する
//             tx := NewTransaction(gormDB)
// 			_ = tx.DoInTx(ctx, func(ctx context.Context) error {
// 				output = greatestUsecase.Do(ctx, tc.input)
//                 if !reflect.DeepEqual(output, GreatestUsecaseOuntput{}) {
//                     t.Fatalf("want: %v, got: %v", tc.want, output)
// 				}
//
//                 // NOTE: テストケースごとに Rollback する
// 				return errors.New("rollback")
// 			})
// 		})
// 	}
// }
//
