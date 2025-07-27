# とりあえずコマンド覚えたいので、使いそうなものは書いていくので、後で直す

tidy:
	go mod tidy

format:
	go fmt ./...

lint:
	go vet ./...

# TODO なんかエラー出る
check:
	staticcheck ./...

# TODO findでファイルを探して、goimportsで整形する
import:
	goimports -w cmd/main.go

unitt:
	go test `go list ./... | grep -v internal/db/query` -v -coverprofile=coverage.out
	# go test -v ./... -coverprofile=coverage.out
	# go tool cover -html=coverage.out -o coverage.html

dbt:
	go test -v ./internal/db/query/... -coverprofile=coverage.out
	# go tool cover -html=coverage.out -o coverage.html

singlet:
	go test -v $(file)

migration:
	migrate create -dir ./scripts/migration -ext sql $(name)
