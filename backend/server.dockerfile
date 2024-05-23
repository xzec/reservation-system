FROM golang:1.22.2

WORKDIR /app
# Download Go modules
COPY go.mod ./
RUN pwd
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY cmd/server/main.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["/server"]