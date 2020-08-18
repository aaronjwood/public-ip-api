FROM golang:1.15-alpine AS builder
WORKDIR /go/src/public-ip-api
COPY go.mod .
COPY server.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/public-ip-api

FROM scratch
COPY --from=builder /go/src/public-ip-api/bin/public-ip-api /public-ip-api
ENTRYPOINT ["/public-ip-api"]
