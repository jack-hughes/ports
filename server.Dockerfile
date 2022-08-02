FROM golang:1.18.5 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o port-domain-service /build/cmd/server/

FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /build/port-domain-service /bin
EXPOSE 50085

CMD ["/bin/port-domain-service"]
