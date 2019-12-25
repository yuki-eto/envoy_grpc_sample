FROM golang:1.13 as builder

WORKDIR /build
COPY main.go .
ENV GO111MODULE=on
RUN go build -o xds_server main.go

FROM debian:buster
COPY --from=builder /build/xds_server .
EXPOSE 30000
CMD ["./xds_server"]
