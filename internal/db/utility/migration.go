package utility

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Migrate(db *sql.DB, pathToRoot string) error {
	// db, err := sql.Open("postgres", "postgres://localhost:5432/database?sslmode=enable")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// 実行rootからの相対pathになっている。全体実行はproject rootだが、テストの場合はtest対象のディレクトリ
	m, err := migrate.NewWithDatabaseInstance("file://"+pathToRoot+"scripts/migration", "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
