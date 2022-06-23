FROM golang:1.17.11 as builder

WORKDIR /code
COPY go.mod go.sum /code/
RUN go mod download

COPY . /code
RUN make build

FROM alpine:3.8
LABEL maintainers="https://tinkerbell.org/community/slack/"

ARG IPMITOOL_REPO=https://github.com/ipmitool/ipmitool.git
ARG IPMITOOL_COMMIT=b5ce925744851b58193ad3ee18957ce88f6efc26
ARG GRPC_HEALTH_PROBE_VERSION=v0.3.4

WORKDIR /tmp
RUN apk add --update --upgrade --no-cache --virtual build-deps \
        alpine-sdk=1.0-r0 \
        autoconf=2.69-r2 \
        automake=1.16.1-r0 \
        git=2.18.4-r0 \
        libtool=2.4.6-r5 \
        ncurses-dev=6.1_p20180818-r1 \
        openssl-dev=1.0.2u-r0 \
        readline-dev=7.0.003-r0 \
    && apk add --update --upgrade --no-cache --virtual run-deps \
	    ca-certificates=20191127-r2 \
        libcrypto1.0=1.0.2u-r0 \
        musl=1.1.19-r11 \
        readline=7.0.003-r0 \
    && git clone -b master ${IPMITOOL_REPO}

WORKDIR /tmp/ipmitool
RUN git checkout ${IPMITOOL_COMMIT} \
    && ./bootstrap \
    && ./configure \
        --prefix=/usr/local \
        --enable-ipmievd \
        --enable-ipmishell \
        --enable-intf-lan \
        --enable-intf-lanplus \
        --enable-intf-open \
    && make \
    && make install \
    && apk del build-deps

WORKDIR /tmp
RUN rm -rf /tmp/ipmitool \
    && wget -O/tmp/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
    && chmod +x /tmp/grpc_health_probe

ENV GIN_MODE release 
USER pbnj
EXPOSE 50051 9090 8080

COPY scripts/etc-passwd /etc/passwd
COPY --from=builder /code/bin/pbnj-linux-amd64 /pbnj-linux-amd64

ENTRYPOINT ["/pbnj-linux-amd64"]
CMD ["server"]
