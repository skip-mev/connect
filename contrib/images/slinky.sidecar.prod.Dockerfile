FROM golang:1.22-bullseye AS builder

WORKDIR /src/slinky
COPY go.mod .

RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/base-debian11:debug
EXPOSE 8080
EXPOSE 8002

COPY --from=builder /src/slinky/build/* /usr/local/bin/

WORKDIR /usr/local/bin/
ENTRYPOINT ["slinky", "--oracle-config-path", "/oracle/oracle.json"]
