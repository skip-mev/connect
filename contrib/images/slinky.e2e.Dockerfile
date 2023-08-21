FROM golang:1.20-bullseye AS builder

WORKDIR /src/slinky
COPY . .

RUN go work init
RUN go work edit -use ./
RUN make tidy
RUN make build-test-app

## Prepare the final clear binary
## This will expose the tendermint and cosmos ports alongside 
## starting up the sim app and the slinky daemon
FROM ubuntu:rolling
EXPOSE 26656 26657 1317 9090 7171
ENTRYPOINT ["slinkyd", "start"]

COPY --from=builder /src/slinky/build/* /usr/local/bin/
RUN apt-get update && apt-get install ca-certificates -y
