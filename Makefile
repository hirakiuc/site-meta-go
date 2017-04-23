.DEFAULT_GOAL := default

build:
	go build

install:
	go install

clean:
	go clean

.PHONY: check
check:
	go vet
	golint

.PHONY: test
test:
	go test -cover

.PHONY: default
default:
	make check
	make test
