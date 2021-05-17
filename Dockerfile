FROM rust:1.52.1-slim as builder
COPY src src
COPY Cargo.toml .
COPY Cargo.lock .
RUN cargo build --release

FROM debian-slim
COPY --from=builder /target/release/public-ip-api /public-ip-api
ENTRYPOINT ["/public-ip-api"]
