FROM ghcr.io/skip-mev/connect-dev-base AS builder

WORKDIR /src/connect

COPY go.mod .

RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/base-debian11:debug
EXPOSE 8080 8002

# Create a non-root user and group
RUN groupadd -g 1337 connect && \
    useradd -u 1337 -g connect -m connect
USER connect

COPY --from=builder /src/connect/build/* /usr/local/bin/

WORKDIR /usr/local/bin/
CMD [ "connect" ]
