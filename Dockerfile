FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
