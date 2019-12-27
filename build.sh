#!/usr/bin/env bash

set -e

function build_image() {
    IMAGE_NAME=$1
    echo "- build $IMAGE_NAME"
    docker build -t ${IMAGE_NAME} .
    docker tag ${IMAGE_NAME}:latest ${ECR_HOST}/${IMAGE_NAME}:latest
    docker push ${ECR_HOST}/${IMAGE_NAME}:latest
}

DOCKER_LOGIN=$(aws ecr get-login --no-include-email --region ap-northeast-1)
eval "${DOCKER_LOGIN}"

cd $(dirname $0)
build_image $APP_IMAGE_NAME

cd ./docker/envoy
build_image $ENVOY_IMAGE_NAME

