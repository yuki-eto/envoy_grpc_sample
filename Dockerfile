FROM toydium/golang-protoc:1.13.5 as builder

RUN apt-get update && apt-get install -y default-mysql-client curl
RUN go get -u github.com/oxequa/realize
RUN go get -u golang.org/x/tools/cmd/goimports

WORKDIR /envoy_grpc_sample
COPY . .
ENV GO111MODULE=on
RUN go build -o app server/main.go

FROM debian:buster
WORKDIR /envoy_grpc_sample
COPY --from=builder /envoy_grpc_sample/app .
#COPY --from=builder /envoy_grpc_sample/bin /envoy_grpc_sample/bin
#COPY ./config /envoy_grpc_sample/config
EXPOSE 30001 30002
CMD ["/envoy_grpc_sample/app"]
