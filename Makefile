lint:
	golangci-lint run

tests:
	go test -v ./... -cover

server:
	go run cmd/server/main.go

proxy:
	go run cmd/reverse/main.go

req:
	curl -i http://localhost:8080