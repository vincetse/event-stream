FROM golang:1.15-alpine AS builder
MAINTAINER vincetse@users.noreply.github.com

RUN \
  apk update && \
  apk upgrade && \
  apk add --upgrade \
    git \
    libprotoc \
    make \
    protobuf-dev \
    protoc \
  && \
  go get google.golang.org/protobuf/cmd/protoc-gen-go && \
  go install google.golang.org/protobuf/cmd/protoc-gen-go

WORKDIR /go/src/github.com/vincetse/event-stream
COPY . ./
RUN \
  make

FROM scratch
COPY --from=builder /go/src/github.com/vincetse/event-stream/event-* /
