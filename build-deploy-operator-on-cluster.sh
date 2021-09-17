#!/bin/bash

IMG=${1:-mustafasaber/custom-k8s-operator:v1.0.0}

docker logout

set -e

current_path=$(pwd)
source "${current_path}"/.docker-secret

echo "$DOCKER_HUB_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin

set +e

make docker-build docker-push IMG="${IMG}"

make deploy IMG="${IMG}"

# kubectl apply -f ./config/samples/_v1_onekind.yaml