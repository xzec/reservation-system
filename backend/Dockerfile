FROM golang:1.22.2 as build

WORKDIR /app

COPY go.mod go.sum /
RUN go mod download

COPY /cmd /cmd
COPY /pkg/server /pkg/server

RUN CGO_ENABLED=0 GOOS=linux go build -C /cmd/server -o /server .

FROM alpine:latest as final

COPY --from=build /server /server

EXPOSE 8081

ENTRYPOINT ["/server"]