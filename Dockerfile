FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/service ./cmd/service

FROM arigaio/atlas@sha256:7159218f66eaed51aa6023c5752f7e87292e536ac6a29eb1194dbf333220784b AS atlas
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/service /app/service
COPY --from=atlas /atlas /usr/local/bin/atlas

COPY configs /app/configs
COPY migrations /app/migrations
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
