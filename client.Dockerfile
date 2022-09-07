FROM golang:1.19.1 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ports-client-api /build/cmd/client/

FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /build/ports-client-api /bin
COPY test/testdata/ports.json /test/ports.json
EXPOSE 50085

CMD ["/bin/ports-client-api"]
