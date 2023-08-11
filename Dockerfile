FROM golang:alpine AS builder

WORKDIR /app

ADD go.mod .

COPY . .

RUN go build -o main ./cmd/app/main.go

# FROM alpine

# WORKDIR /app

# COPY --from=builder /app/main /app/main

CMD ["./main"]