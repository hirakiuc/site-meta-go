.DEFAULT_GOAL := default

build:
	go build

install:
	go install

clean:
	go clean

check:
	go vet
	golint

test:
	go test -cover

.PHONY: default
default:
	make check
	make test
