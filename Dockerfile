FROM golang:1.15 as builder
LABEL maintainers="https://tinkerbell.org/community/slack/"

ARG GRPC_HEALTH_PROBE_VERSION=v0.3.4

RUN wget -O/tmp/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
		chmod +x /tmp/grpc_health_probe

WORKDIR /code
COPY go.mod go.sum /code/
RUN go mod download

COPY . /code
RUN make build

FROM scratch
USER pbnj
EXPOSE 9090

COPY scripts/etc-passwd /etc/passwd
COPY --from=builder /code/bin/pbnj-linux-amd64 /pbnj-linux-amd64
COPY --from=builder /tmp/grpc_health_probe /bin/grpc_health_probe

ENTRYPOINT ["/pbnj-linux-amd64"]
CMD ["server"]
