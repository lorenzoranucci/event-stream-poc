FROM alpine:latest as certs

RUN apk --update add ca-certificates

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./bin/golang-kafka /bin/golang-kafka

EXPOSE 8080

ENTRYPOINT ["/bin/golang-kafka"]
