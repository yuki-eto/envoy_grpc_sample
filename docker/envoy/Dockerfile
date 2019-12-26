FROM envoyproxy/envoy:latest

COPY ./start_envoy.sh .
COPY ./conf/*.yaml /etc/envoy/
EXPOSE 8001 8002 9000
CMD ./start_envoy.sh $ENVOY_ENV
