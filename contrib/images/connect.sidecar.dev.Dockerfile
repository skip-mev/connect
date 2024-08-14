FROM ghcr.io/skip-mev/connect-dev-base AS builder

WORKDIR /src/connect

COPY go.mod .

RUN go mod download

COPY . .

RUN make build

FROM ubuntu:rolling
EXPOSE 8080 8002

COPY --from=builder /src/connect/build/* /usr/local/bin/
RUN apt-get update && apt-get install jq -y && apt-get install ca-certificates -y

WORKDIR /usr/local/bin/
ENTRYPOINT [ "connect" ]
