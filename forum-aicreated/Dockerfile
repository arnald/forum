FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o forum ./cmd/forum

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/

COPY --from=builder /app/forum .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./forum"]