.PHONY: build
build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/public-api-server

.PHONY: image
image:
	@docker build -t public-ip-api .
