# Step 1: Modules caching
FROM golang:1.21-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.21-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/bot ./cmd/main.go

# Step 3: Final
FROM alpine:latest
RUN apk --no-cache add ca-certificates

EXPOSE 8001

# GOPATH for scratch images is /
COPY --from=builder /app/.env /
COPY --from=builder /bin/bot /bot
CMD ["/bot"]