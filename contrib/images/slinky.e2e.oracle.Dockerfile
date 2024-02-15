FROM golang:1.21-bullseye AS builder

WORKDIR /src/slinky
COPY go.mod .

RUN go mod download

COPY . .

RUN make build

FROM ubuntu:rolling
EXPOSE 8080
EXPOSE 8001

COPY --from=builder /src/slinky/build/* /usr/local/bin/
RUN apt-get update && apt-get install ca-certificates -y

WORKDIR /usr/local/bin/
ENTRYPOINT ["oracle", "--oracle-config-path", "/oracle/oracle.json", "--market-config-path", "/oracle/market.json", "-host", "0.0.0.0", "-port", "8080"]
