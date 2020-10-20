FROM alpine:3.8

ARG IPMITOOL_REPO=https://github.com/ipmitool/ipmitool.git
ARG IPMITOOL_COMMIT=c3939dac2c060651361fc71516806f9ab8c38901

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

CMD ["/pbnj"]
ENV GIN_MODE release
EXPOSE 9090
RUN adduser -h /home/packet -s /bin/sh -D packet
USER packet
COPY pbnj /
