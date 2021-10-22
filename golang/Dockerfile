FROM alpine:latest as certs

RUN apk --update add ca-certificates

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./bin/event-stream-golang /bin/event-stream-golang

EXPOSE 8080

ENTRYPOINT ["/bin/event-stream-golang"]
