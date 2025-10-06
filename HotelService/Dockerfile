FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o hotel-service ./server/cmd/hotel

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/hotel-service .

ENTRYPOINT ["/app/hotel-service"]