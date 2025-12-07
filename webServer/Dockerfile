FROM golang:1.25.5 AS builder
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o webServer main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /webServer .
CMD ["./webServer"]