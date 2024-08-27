FROM ghcr.io/skip-mev/connect-dev-base AS builder

WORKDIR /src/connect

COPY go.mod .

RUN go mod download

COPY . .

RUN make build

# Create a non-root user and group
RUN groupadd -r connectgroup && useradd -r -g connectgroup connect

FROM gcr.io/distroless/base-debian11:debug
EXPOSE 8080 8002

# Copy the user and group files from the builder stage
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Set the user to connect
USER connect

COPY --from=builder /src/connect/build/* /usr/local/bin/
COPY --from=builder /src/connect/scripts/deprecated-exec.sh /usr/local/bin/

WORKDIR /usr/local/bin/
ENTRYPOINT [ "deprecated-exec.sh" ]
CMD [ "slinky" ]
