build:
	GOOS=linux GOARCH=arm64 go build -o cmd/webhook/bootstrap cmd/webhook/main.go
	zip -j webhook.zip cmd/webhook/bootstrap

	GOOS=linux GOARCH=arm64 go build -o cmd/worker/bootstrap cmd/worker/main.go
	zip -j worker.zip cmd/worker/bootstrap
