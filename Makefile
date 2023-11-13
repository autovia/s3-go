BINARY_NAME=s3-go
OS=linux

build:
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-${OS}-amd64 main.go
	GOARCH=arm64 GOOS=linux go build -o bin/${BINARY_NAME}-${OS}-arm64 main.go