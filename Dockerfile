FROM rust:1.52.1-slim as builder
COPY src .
RUN cargo build --release

FROM scratch
COPY --from=builder /target/release/public-ip-api /public-ip-api
ENTRYPOINT ["/public-ip-api"]
