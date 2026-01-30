FROM alpine:latest as builder

RUN apk --no-cache add ca-certificates

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD release/linux/amd64/drone-env-inspector /bin/

ENTRYPOINT ["/bin/drone-env-inspector"]

