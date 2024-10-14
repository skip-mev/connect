FROM ghcr.io/skip-mev/connect-dev-base AS builder

WORKDIR /src/connect

COPY go.mod .

RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/base-debian11:debug
EXPOSE 8080 8002

RUN ["/busybox/busybox","adduser","-u","1000","-g","1000","-D","--home","/home/connect","-s","/sbin/nologin","connect"]

COPY --from=builder /src/connect/build/* /usr/local/bin/

WORKDIR /usr/local/bin/
USER connect
CMD [ "connect" ]
