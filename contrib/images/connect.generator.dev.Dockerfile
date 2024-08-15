# ./contrib/images/connect.generator.dev.Dockerfile

# Stage 1: Build the Go application
FROM golang:1.22 AS builder

WORKDIR /src/connect

COPY go.mod .

RUN go mod download

COPY . .

RUN make build

# Stage 2: Create a lightweight image for running the application
FROM ubuntu:rolling
COPY --from=builder /src/connect/build/* /usr/local/bin/

# Create the /data directory
RUN mkdir -p /data
# Define the volume where the generated file will be stored
VOLUME /data

# The entrypoint will be provided by the docker-compose file
ENTRYPOINT ["/usr/local/bin/scripts"]
