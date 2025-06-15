FROM golang:1.21 AS builder

WORKDIR /build
COPY . .
RUN go mod init dude && go mod tidy && go build -o app

FROM golang:1.21

WORKDIR /app
COPY --from=builder /build /app
CMD ["./app"]
