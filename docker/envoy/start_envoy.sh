#!/usr/bin/env bash

ENVOY_ENV=$1
if [[ "$ENVOY_ENV" = "" ]]; then
    ENVOY_ENV="local"
fi

ENVOY_CONF_NAME="front-envoy-$ENVOY_ENV.yaml"
envoy -c /etc/envoy/${ENVOY_CONF_NAME}
