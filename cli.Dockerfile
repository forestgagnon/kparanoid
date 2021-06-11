FROM golang:1.16-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -v -ldflags="-s -w" -o kparanoid ./bin/cli

FROM alpine:3.13
RUN apk add --update --no-cache ca-certificates curl docker-cli
WORKDIR /app/

COPY installation installation
COPY --from=builder /app/kparanoid ./
ENTRYPOINT ["/app/kparanoid"]
