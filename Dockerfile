# 1st stage: develop environment
FROM toydium/golang-protoc:1.13.5 as develop
WORKDIR /envoy_grpc_sample
RUN apt-get update && apt-get install -y default-mysql-client curl
ENV GO111MODULE=off
RUN go get -u github.com/oxequa/realize
RUN go get -u golang.org/x/tools/cmd/goimports
ENV GO111MODULE=on
COPY . .

# 2nd stage: build environment
FROM golang:1.13.5 as builder
WORKDIR /envoy_grpc_sample
COPY . .
ENV GO111MODULE=on
ENV GOPROXY=http://127.0.0.1:8009
RUN go build -o app server/main.go

# 3rd stage: running environment
FROM debian:buster
WORKDIR /envoy_grpc_sample
COPY --from=builder /envoy_grpc_sample/app .
#COPY --from=builder /envoy_grpc_sample/bin /envoy_grpc_sample/bin
#COPY ./config /envoy_grpc_sample/config
EXPOSE 30001 30002
CMD ["/envoy_grpc_sample/app"]
