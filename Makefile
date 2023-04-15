lint:
	golangci-lint run

tests:
	go test -v ./... -cover

req:
	curl -i http://localhost:8080/api/v1/health?name=foo

docker:
	docker-compose --file build/docker-compose.yml up --build