#!/usr/bin/env bash

cd $(dirname "$0") || exit

yae ./front-envoy-local-base.yaml > ./conf/front-envoy-local.yaml
yae ./front-envoy-ecs-base.yaml > ./conf/front-envoy-ecs.yaml
