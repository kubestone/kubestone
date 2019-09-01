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
	samplesDir = "../../config/samples/"
)

var restClientConfig = ctrl.GetConfigOrDie()
var client ctrlclient.Client
var ctx = context.Background()
var scheme = runtime.NewScheme()

func init() {
	_ = k8sscheme.AddToScheme(scheme)
	_ = perfv1alpha1.AddToScheme(scheme)

	var err error
	client, err = ctrlclient.New(restClientConfig, ctrlclient.Options{Scheme: scheme})
	if err != nil {
		panic(err)
	}
}

func TestEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "end-to-end suite")
}
