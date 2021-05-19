.PHONY: build
build:
	@cargo build --release

.PHONY: image
image:
	@docker build -t public-ip-api .
