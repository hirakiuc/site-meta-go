FROM golang:1.14-alpine AS build
WORKDIR /go/src/site-meta-go
COPY . .
RUN apk add --no-cache --virtual dev-deps make git gcc libc-dev \
  && go get -d -v ./... \
  && go mod download \
  && go mod vendor \
  && make build-server \
  && apk del dev-deps

FROM golang:1.14-alpine
COPY --from=build /go/src/site-meta-go/api-server /bin/server
EXPOSE 8080
ENTRYPOINT ["/bin/server"]
