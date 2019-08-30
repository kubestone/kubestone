#!/bin/bash -x

# Copyright 2019 The xridge kubestone contributors.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eEuo pipefail

DOCKER_IMAGE="xridge/kubestone:e2e"
KUBESTONE_ROOT=$(dirname $0)/../../../

build_kubestone() {
    pushd ${KUBESTONE_ROOT}
    docker build -t ${DOCKER_IMAGE} .
    popd
}

upload_kubestone_to_kind() {
    kind load --loglevel debug docker-image ${DOCKER_IMAGE}
}

deploy_kubestone() {
    pushd ${KUBESTONE_ROOT}
    make deploy-e2e
    popd
}

validate_kubestone_deployment() {
    kubectl -n kubestone-system \
        wait --for=condition=Available --timeout=1m \
        deployments/kubestone-controller-manager
}

show_all_objects() {
    kubectl get all --all-namespaces
}

main() {
    build_kubestone
    upload_kubestone_to_kind
    deploy_kubestone
    validate_kubestone_deployment
    show_all_objects
}

main

