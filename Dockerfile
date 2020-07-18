FROM scratch
COPY bin/public-api-server /public-api-server
ENTRYPOINT ["/public-api-server"]