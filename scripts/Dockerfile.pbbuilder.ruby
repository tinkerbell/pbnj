FROM ruby:3
RUN apt update
RUN apt install -y golang
RUN go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@v0.3.2
RUN gem install grpc-tools -v 1.27
