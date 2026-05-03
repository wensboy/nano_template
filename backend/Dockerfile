FROM golang:1.25.9-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o main .
RUN apk --no-cache add ca-certificates

FROM scratch
COPY ./.env ./env.yaml /
COPY --from=builder /app/main /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 3000

# maybe some start flags here...
CMD ["/main"]