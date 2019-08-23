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


# This is a modified version of kind deploy script, taken from:
# https://github.com/kind-ci/examples

set -eEuo pipefail

BIN_DIR="$(mktemp -d)"
KIND="${BIN_DIR}/kind"
KUBECTL="${BIN_DIR}/kubectl"

cleanup() {
    kubectl get all --all-namespaces || true
    # TODO: Lack of labels on deployment makes it impossible to use selectors
    kubectl logs -n kubestone-system $(kubectl get pods -n kubestone-system -o name | grep controller) || true
    "${KIND}" delete cluster || true
    rm -rf "${BIN_DIR}"
}
trap cleanup EXIT

OS=$(uname -s | tr A-Z a-z)

install_latest_kind() {
    # clone kind into a tempdir within BIN_DIR
    local tmp_dir
    tmp_dir="$(TMPDIR="${BIN_DIR}" mktemp -d "${BIN_DIR}/kind-source.XXXXX")"
    cd "${tmp_dir}" || exit
    git clone https://github.com/kubernetes-sigs/kind && cd ./kind
    make install INSTALL_DIR="${BIN_DIR}"
}

install_kind_release() {
    KIND_VERSION="v0.4.0"
    KIND_BINARY_URL="https://github.com/kubernetes-sigs/kind/releases/download/${KIND_VERSION}/kind-${OS}-amd64"
    curl -L -o "${KIND}" "${KIND_BINARY_URL}"
    chmod +x "${KIND}"
}

install_kubectl_release() {
    KUBECTL_VERSION="v1.15.2"
    KUBECTL_BINARY_URL="https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/${OS}/amd64/kubectl"

    curl -L -o "${KUBECTL}" "${KUBECTL_BINARY_URL}"
    chmod +x "${KUBECTL}"
}

main() {
    export PATH=${BIN_DIR}:${PATH}

    install_kind_release
    kind_config="$(dirname $0)/../conf/kind-two-worker-cluster.yaml"
    "${KIND}" create cluster --loglevel=debug --config ${kind_config}
    KUBECONFIG="$("${KIND}" get kubeconfig-path)"
    export KUBECONFIG

    install_kubectl_release
    kubectl version

    $(dirname $0)/deploy-kubestone.sh

    go test -v ./tests/e2e/...
}

main
