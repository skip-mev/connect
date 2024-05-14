FROM ghcr.io/skip-mev/slinky-dev-base as builder

WORKDIR /src/slinky

COPY go.mod .

RUN go mod download

COPY . .

RUN make build

FROM ubuntu:rolling
EXPOSE 8080 8002

COPY --from=builder /src/slinky/build/* /usr/local/bin/
RUN apt-get update && apt-get install jq -y && apt-get install ca-certificates -y

WORKDIR /usr/local/bin/
ENTRYPOINT ["slinky", "--oracle-config-path", "/oracle/oracle.json", "--update-market-config-path", "/oracle/market.json", "--market-map-endpoint", "blockchain:9090"]