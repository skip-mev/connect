FROM golang:1.21-bullseye AS builder

WORKDIR /src/slinky
COPY . .
RUN make tidy
RUN make build

FROM ubuntu:rolling
EXPOSE 8080


COPY --from=builder /src/slinky/build/* /usr/local/bin/
RUN apt-get update && apt-get install ca-certificates -y

WORKDIR /usr/local/bin/
ENTRYPOINT ["oracle", "-config", "/oracle/config.toml", "-host", "0.0.0.0", "-port", "8080"]
