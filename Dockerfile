FROM public.ecr.aws/docker/library/golang:1.21-alpine AS builder
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
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

WORKDIR /go/src/github.com/vincetse/event-stream
COPY . ./
RUN \
  make

FROM scratch
COPY --from=builder /go/src/github.com/vincetse/event-stream/event-* /
