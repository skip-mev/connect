FROM ghcr.io/skip-mev/connect-dev-base AS builder

WORKDIR /src/connect

COPY go.mod .

RUN go mod download

COPY . .

RUN make build-test-app

## Prepare the final clear binary
## This will expose the tendermint and cosmos ports alongside 
## starting up the sim app and the connect daemon
FROM ubuntu:rolling
EXPOSE 26656 26657 1317 9090 7171 26655 8081 26660

RUN apt-get update && apt-get install jq -y && apt-get install ca-certificates -y
ENTRYPOINT ["connectd", "start"]

COPY --from=builder /src/connect/build/* /usr/local/bin/
