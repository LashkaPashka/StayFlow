FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o payment-service ./server/cmd/payment

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/payment-service .

ENTRYPOINT ["/app/payment-service"]