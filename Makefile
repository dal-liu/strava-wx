build:
	GOOS=linux GOARCH=arm64 go build -o bootstrap main.go
	zip lambda_handler.zip bootstrap
