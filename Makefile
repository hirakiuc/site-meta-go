.DEFAULT_GOAL := default

SERVER=api-server

.PHONY: build
build: build-server

.PHONY:build-server
build-server:
	go build -o ${SERVER} ./cmd/server/main.go

.PHONY: run-server
run-server:
	go run ./cmd/server/main.go

.PHONY: clean
clean:
	go clean
	rm -f ${SERVER}

.PHONY: check
check:
	golangci-lint run --enable-all ./...

.PHONY: test
test:
	go test -cover ./...

.PHONY: testbuild
testbuild:
	go test -c -args -w -gcflags "-N -l" ./...

.PHONY: default
default:
	make check
	make test
