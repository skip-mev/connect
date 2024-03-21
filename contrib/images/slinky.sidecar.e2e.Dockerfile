FROM golang:1.22-bullseye AS builder

WORKDIR /src/slinky
COPY go.mod .

RUN go mod download

COPY . .

RUN make build
RUN make update-local-configs


FROM ubuntu:rolling
EXPOSE 8080
EXPOSE 8002

COPY --from=builder /src/slinky/build/* /usr/local/bin/
COPY --from=builder /src/slinky/config/local /etc/slinky/default_config
RUN apt-get update && apt-get install ca-certificates -y

WORKDIR /usr/local/bin/
ENTRYPOINT ["oracle", "--oracle-config-path", "/oracle/oracle.json", "--market-config-path", "/oracle/market.json"]
