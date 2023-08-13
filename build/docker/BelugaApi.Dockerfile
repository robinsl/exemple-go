FROM golang:1.20 as builder
LABEL stage="builder"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod tidy && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -v -o beluga-api ./cmd/beluga-api

CMD ["ls", "/app"]

FROM scratch as production
LABEL stage="production"

WORKDIR /app

COPY --from=builder /app/beluga-api .
COPY --from=builder /app/configs/local.env ./configs/.env

EXPOSE 80

CMD ["./beluga-api"]

