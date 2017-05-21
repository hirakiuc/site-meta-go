.DEFAULT_GOAL := default

build:
	go build

install:
	go install

clean:
	go clean

.PHONY: check
check:
	go vet . ./internal/...
	golint ./internal/...
	golint .

.PHONY: test
test:
	go test -cover

testbuild:
	go test -c -args -w -gcflags "-N -l"

.PHONY: default
default:
	make check
	make test
