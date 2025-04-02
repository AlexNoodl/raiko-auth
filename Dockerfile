FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o raiko-auth ./cmd

FROM alpine:latest
#RUN apk add --no-cache ca-certificates
WORKDIR /app/

COPY --from=builder /app/raiko-auth /app/raiko-auth

COPY --from=builder /app/.env /app/.env

RUN chmod +x /app/raiko-auth

EXPOSE 8080

CMD ["/app/raiko-auth"]