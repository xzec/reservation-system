FROM golang:1.22.2 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/server/main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

EXPOSE 8080

FROM alpine:latest as final

WORKDIR /

COPY --from=build /server /server

EXPOSE 8080

ENTRYPOINT ["/server"]