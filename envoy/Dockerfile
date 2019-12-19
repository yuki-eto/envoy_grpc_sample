FROM envoyproxy/envoy:latest

COPY ./front-envoy-ecs.yaml /etc/envoy/envoy.yaml
EXPOSE 8001 8002 9000
CMD ["envoy", "-c", "/etc/envoy/envoy.yaml"]
