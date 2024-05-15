FROM ghcr.io/skip-mev/slinky-dev-base as builder

WORKDIR /src/slinky

COPY go.mod .

RUN go mod download

COPY . .

RUN make build-test-app

## Prepare the final clear binary
## This will expose the tendermint and cosmos ports alongside 
## starting up the sim app and the slinky daemon
EXPOSE 26656 26657 1317 9090 7171 26655 8081 26660
RUN apt-get update && apt-get install jq -y && apt-get install ca-certificates -y
ENTRYPOINT ["make", "build-and-start-app"]

