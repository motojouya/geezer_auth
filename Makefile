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
