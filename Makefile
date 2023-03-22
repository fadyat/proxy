lint:
	golangci-lint run

tests:
	go test -v ./... -cover

server:
	go run cmd/server/main.go

reverse:
	go run cmd/reverse/main.go

request:
	curl -i http://localhost:8080