lint:
	golangci-lint run

tests:
	go test -v ./... -cover

req:
	curl -i http://localhost:8080/api/v1/health?name=foo2 \
		-H "Content-Type: application/json" \
		-X 'GET' \
		--data '{"name":"foo"}'

docker:
	docker-compose --file build/docker-compose.yml up --build proxy

cache:
	docker-compose --file build/docker-compose.yml up --build cache-local
