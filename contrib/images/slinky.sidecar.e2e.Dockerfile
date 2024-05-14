FROM ghcr.io/skip-mev/slinky-dev-base as builder

WORKDIR /src/slinky

COPY go.mod .

RUN go mod download

COPY . .

RUN make build
RUN make update-local-configs

FROM ubuntu:rolling
EXPOSE 8080 8002

COPY --from=builder /src/slinky/build/* /usr/local/bin/
COPY --from=builder /src/slinky/config/local /etc/slinky/default_config
RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /usr/local/bin/
ENTRYPOINT ["slinky", "--oracle-config-path", "/oracle/oracle.json", "--update-market-config-path", "/oracle/market.json"]
