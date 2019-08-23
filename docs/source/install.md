# Installation

## Prerequisites
* [Kubernetes](https://kubernetes.io) v1.14 (or newer)
* [Kustomize](https://kustomize.io) v3.10

## Deploy the Operator

The following command will

* create namespace `kubestone-system` for the kubestone operator
* install the CRDs to the cluster
* create the necessary roles and rolebindings
* install and start the operator

```bash
$ kustomize build config/default | kubectl apply -f -
namespace/kubestone-system created
customresourcedefinition.apiextensions.k8s.io/fios.perf.kubestone.xridge.io configured
customresourcedefinition.apiextensions.k8s.io/iperf3s.perf.kubestone.xridge.io configured
role.rbac.authorization.k8s.io/kubestone-leader-election-role created
clusterrole.rbac.authorization.k8s.io/kubestone-manager-role created
clusterrole.rbac.authorization.k8s.io/kubestone-proxy-role created
rolebinding.rbac.authorization.k8s.io/kubestone-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/kubestone-manager-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/kubestone-proxy-rolebinding created
service/kubestone-controller-manager-metrics-service created
deployment.apps/kubestone-controller-manager created
```
