FROM golang:1.19-alpine as builder
RUN apk update && apk add --no-cache git make ca-certificates tzdata && update-ca-certificates
RUN mkdir /subscriptions
WORKDIR /subscriptions
COPY . .
WORKDIR /subscriptions
RUN make build

FROM alpine
ENV USER=runner USER_ID=1002 USER_G=runner USER_G_ID=1002
RUN addgroup -g ${USER_G_ID} ${USER_G} && \
    adduser -D --home /app -u ${USER_ID} -G ${USER_G} ${USER}
WORKDIR /app
COPY --from=builder /subscriptions/bin/subscriptions /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/subscriptions"]
