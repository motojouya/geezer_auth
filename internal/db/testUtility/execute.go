package testUtility

import (
	//_ "internal/core/timezone"
	"database/sql"
	"fmt"
	"github.com/motojouya/geezer_auth/internal/db"
	"github.com/motojouya/geezer_auth/internal/db/utility"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
)

func ExecuteDatabaseTest(pathToRoot string, run func(db.ORP) int) {
	os.Setenv("TZ", "Asia/Tokyo") // `internal/core/timezone`だとできんかった？

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=geezer_auth",
			"POSTGRES_USER=geezer_auth",
			"POSTGRES_DB=geezer_auth",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// resource.Expire(120) // Tell docker to hard kill the container in 120 seconds
	// pool.MaxWait = 120 * time.Second // exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	var database *sql.DB
	databaseUrl := fmt.Sprintf("postgres://geezer_auth:geezer_auth@%s/geezer_auth?sslmode=disable", resource.GetHostPort("5432/tcp"))
	if err = pool.Retry(func() error {
		database, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return database.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	var migrateErr = utility.Migrate(database, pathToRoot)
	if migrateErr != nil {
		log.Fatalf("Could not migrate database: %s", migrateErr)
	}

	var orp = db.CreateDatabase(database)

	code := run(orp)

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
