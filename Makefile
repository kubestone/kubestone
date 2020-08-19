
# Image URL to use all building/pushing image targets
IMG ?= kubestone:latest

API_VERSION ?= "v1alpha1"

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the last release tag if an override is not provided
KUBESTONE_RELEASE ?= $(shell git tag -l | egrep "v\d+\.\d+\.\d+" | tail -1)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOLANGCI_VERSION = v1.21.0

all: manager

# Run unit tests
test: generate fmt lint manifests
	go test -v ./api/... ./controllers/... ./pkg/... -coverprofile cover.out

# Run end to end tests
e2e-test:
	./tests/e2e/bin/e2e-test-in-kind.sh

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd/bases

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kustomize build config/default | kubectl apply -f -

# Deployment used for end-to-end test
deploy-e2e: manifests
	kustomize build config/e2e | kubectl apply -f -

# Generate manifests: CRD, RBAC, etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Download gen-crd-api-reference-docs
gen-crd-api-reference-docs:
ifeq (, $(shell which gen-crd-api-reference-docs))
	go get github.com/ahmetb/gen-crd-api-reference-docs@v0.1.5
GEN_CRD_API_REFERENCE_DOCS=$(shell go env GOPATH)/bin/gen-crd-api-reference-docs
else
GEN_CRD_API_REFERENCE_DOCS=$(shell which gen-crd-api-reference-docs)
endif

# Generate apidocs
apidocs: gen-crd-api-reference-docs
	$(GEN_CRD_API_REFERENCE_DOCS) -config docs/apidocs/config.json -api-dir github.com/xridge/kubestone/api/$(API_VERSION)/ -out-file docs/static/apidocs.html -template-dir docs/apidocs/template/

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Download golangci-lint if needed
golangci-lint:
ifeq (, $(shell which golangci-lint))
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_VERSION}
GOLANGCI_LINT=$(shell go env GOPATH)/bin/golangci-lint
else
GOLANGCI_LINT=$(shell which golangci-lint)
endif

# Run linter. GOGC is set to reduce memory footprint
lint: golangci-lint
	GOGC=10 $(GOLANGCI_LINT) run -v --timeout 10m

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./api/...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

# All the things needed before you make a PR
pre-commit: generate apidocs manifests fmt vet
	@echo "Updating quickstart doc to current release ${KUBESTONE_RELEASE}"
	sed -i'' -E "s@(github\.com/xridge/kubestone/config/default)(\?ref=v[0-9]+\.[0-9]+\.[0-9]+)+\b@\1?ref=${KUBESTONE_RELEASE}@g" docs/quickstart.md
