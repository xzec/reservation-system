FROM golang:1.22.2 as build

WORKDIR /app

COPY go.mod go.sum /
RUN go mod download

COPY /cmd/migrate /cmd/migrate
COPY /pkg/db/migrations /pkg/db/migrations

RUN CGO_ENABLED=0 GOOS=linux go build -C /cmd/migrate -o /migrate .

FROM alpine:latest as final

COPY --from=build /migrate /migrate
COPY --from=build /pkg/db/migrations /migrations

CMD ["./migrate", "-dir=migrations", "up"]
