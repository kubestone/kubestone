/*
Copyright 2019 The xridge kubestone contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"log"
	"testing"

	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/runtime"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

const (
	samplesDir = "../../config/samples"
	testConf   = "./conf"
)

const (
	e2eNamespaceDrill    = "kubestone-e2e-drill"
	e2eNamespaceFio      = "kubestone-e2e-fio"
	e2eNamespaceIoping   = "kubestone-e2e-ioping"
	e2eNamespaceIperf3   = "kubestone-e2e-iperf3"
	e2eNamespacePgbench  = "kubestone-e2e-pgbench"
	e2eNamespaceQperf    = "kubestone-e2e-qperf"
	e2eNamespaceSysbench = "kubestone-e2e-sysbench"
)

var e2eNamespaces = []string{
	e2eNamespaceDrill,
	e2eNamespaceFio,
	e2eNamespaceIoping,
	e2eNamespaceIperf3,
	e2eNamespacePgbench,
	e2eNamespaceQperf,
	e2eNamespaceSysbench,
}

var restClientConfig = ctrl.GetConfigOrDie()
var client ctrlclient.Client
var ctx = context.Background()
var scheme = runtime.NewScheme()

var _ = BeforeSuite(func() {
	_ = k8sscheme.AddToScheme(scheme)
	_ = perfv1alpha1.AddToScheme(scheme)

	var err error
	client, err = ctrlclient.New(restClientConfig, ctrlclient.Options{Scheme: scheme})
	if err != nil {
		Fail(err.Error())
	}

	for _, namespace := range e2eNamespaces {
		_, _, err = run("kubectl create namespace " + namespace)
		if err != nil {
			Fail(err.Error())
		}
	}
})

var _ = AfterSuite(func() {
	for _, namespace := range e2eNamespaces {
		stdout, _, err := run("kubectl get all --namespace " + namespace)
		if err != nil {
			Fail(err.Error())
		}
		log.Printf("objects in %s namespace:\n%s\n", namespace, stdout)

		_, _, err = run("kubectl delete namespace " + namespace)
		if err != nil {
			Fail(err.Error())
		}
	}
})

func TestEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "end-to-end suite")
}
