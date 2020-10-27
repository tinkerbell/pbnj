FROM golang:1.15 as builder
LABEL maintainers="https://tinkerbell.org/community/slack/"

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

ENTRYPOINT ["/pbnj-linux-amd64"]
CMD ["server"]
