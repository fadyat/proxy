FROM golang:1.20-alpine3.17 AS builder
WORKDIR /code

ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go-build \
    --mount=type=cache,target=/go/pkg/mod/cache \
    go build -o bin/proxy cmd/reverse/main.go

FROM scratch
COPY --from=builder /code/bin/proxy /usr/local/bin/proxy
ENTRYPOINT ["/usr/local/bin/proxy"]