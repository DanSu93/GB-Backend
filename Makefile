build:
	go build ./...

check:
	golangci-lint run

test:
	go test ./...