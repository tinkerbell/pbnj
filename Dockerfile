FROM golang:1.15 as builder

WORKDIR /code
COPY go.mod go.sum /code/
RUN go mod download

COPY . /code
RUN make build

FROM alpine:3.8
LABEL maintainers="https://tinkerbell.org/community/slack/"

ARG IPMITOOL_REPO=https://github.com/ipmitool/ipmitool.git
ARG IPMITOOL_COMMIT=c3939dac2c060651361fc71516806f9ab8c38901
ARG GRPC_HEALTH_PROBE_VERSION=v0.3.4

WORKDIR /tmp
RUN apk add --update --upgrade --no-cache --virtual build-deps \
        alpine-sdk \
        autoconf \
        automake \
        git \
        libtool \
        ncurses-dev \
        openssl-dev \
        readline-dev \
        && \
    apk add --update --upgrade --no-cache --virtual pbnj-runtime-deps \
	ca-certificates \
        libcrypto1.0 \
        musl \
        readline \
        && \
    git clone -b master ${IPMITOOL_REPO} && \
    cd ipmitool && \
    git checkout ${IPMITOOL_COMMIT} && \
    ./bootstrap && \
    ./configure \
        --prefix=/usr/local \
        --enable-ipmievd \
        --enable-ipmishell \
        --enable-intf-lan \
        --enable-intf-lanplus \
        --enable-intf-open \
        && \
    make && \
    make install && \
    cd $OLDPWD && \
    rm -rf /tmp/ipmitool && \
    apk del build-deps

RUN wget -O/tmp/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
		chmod +x /tmp/grpc_health_probe

ENV GIN_MODE release 
USER pbnj
EXPOSE 50051 9090

COPY scripts/etc-passwd /etc/passwd
COPY --from=builder /code/bin/pbnj-linux-amd64 /pbnj-linux-amd64

ENTRYPOINT ["/pbnj-linux-amd64"]
CMD ["server"]
