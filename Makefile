.PHONY: build
build:
	@cargo build --release

.PHONY: image
image:
	@docker build --pull -t public-ip-api .
