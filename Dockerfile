FROM golang:1.23.0-bookworm as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./cmd/main.go

FROM gcr.io/distroless/base-debian11

COPY --from=builder /bot /bot

ENTRYPOINT ["/bot"]