export GO_PATH := /tmp/go
export PATH := $(GO_PATH)/bin:$(PATH)

PROTOC_OUT = ./build/gen
go = CGO_ENABLED=0 go

all: gen deps build

build:
	$(go) build ./cmd/event-producer/...
	$(go) build ./cmd/event-consumer/...

gen:
	protoc --go_out=. --go_opt=paths=source_relative ./pkg/event/v1/event.proto

docker:
	docker build --rm -t data-stream .

lint:
	test -z `gofmt -s -l .`
	go vet ./...
	golint -set_exit_status `go list ./...`

deps:
	#go get ./...

mq:
	docker run --rm \
		-p 5672:5672 \
		--detach \
		--name mq \
		rabbitmq:3-alpine
