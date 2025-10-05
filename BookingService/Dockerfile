FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o booking-service ./cmd/booking

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/booking-service .

ENTRYPOINT ["/app/booking-service"]