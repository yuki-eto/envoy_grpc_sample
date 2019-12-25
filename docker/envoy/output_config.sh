#!/usr/bin/env bash

cd $(dirname "$0") || exit

yae ./front-envoy-local-base.yaml > ./front-envoy-local.yaml
yae ./front-envoy-ecs-base.yaml > ./front-envoy-ecs.yaml
