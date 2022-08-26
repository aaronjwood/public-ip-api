FROM rust:1-slim-bullseye as builder
WORKDIR /usr/src
RUN USER=root cargo new public-ip-api
COPY Cargo.toml Cargo.lock /usr/src/public-ip-api/
WORKDIR /usr/src/public-ip-api
RUN cargo build --release

FROM scratch
COPY --from=builder /usr/src/public-ip-api/target/release/public-ip-api /public-ip-api
ENTRYPOINT ["/public-ip-api"]
