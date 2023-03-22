lint:
	golangci-lint run

tests:
	go test -v ./... -cover