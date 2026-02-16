# STEP 1 build executable binary
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache tzdata
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy dependency modules first (for caching)
COPY orm/ ./orm/
COPY ez-client-go/ ./ez-client-go/

# Copy go mod and sum files
COPY pocket-shop/go.mod pocket-shop/go.sum ./pocket-shop/

WORKDIR /app/pocket-shop
RUN go mod download

# Copy the source code
COPY pocket-shop/ ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -tags=sonic -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/app cmd/api/main.go

RUN chmod 0555 /go/bin/app

# STEP 2 build a small image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

USER appuser

WORKDIR /app

COPY --from=builder /go/bin/app ./app_binary
COPY --from=builder /app/pocket-shop/docs ./docs

EXPOSE 8080

CMD ["./app_binary"]
