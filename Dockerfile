FROM golang:1.23rc2-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -ldflags "-s -w" -o /app/ds389-exporter ds389-exporter.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/ds389-exporter .

ENTRYPOINT ["/root/ds389-exporter"]

CMD []
