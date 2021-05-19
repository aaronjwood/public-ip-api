FROM rust:1.52.1-slim as builder
RUN apt-get update && apt-get install -y musl-tools
RUN rustup target add x86_64-unknown-linux-musl
WORKDIR /usr/src
RUN USER=root cargo new public-ip-api
COPY Cargo.toml Cargo.lock /usr/src/public-ip-api/
WORKDIR /usr/src/public-ip-api
RUN cargo build --release
COPY src /usr/src/public-ip-api/src
RUN cargo install --target x86_64-unknown-linux-musl --path .

FROM scratch
COPY --from=builder /usr/local/cargo/bin/public-ip-api /public-ip-api
ENTRYPOINT ["/public-ip-api"]
