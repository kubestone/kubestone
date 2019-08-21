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

package main

import (
	"flag"
	"os"

	"github.com/go-logr/zapr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
	"github.com/xridge/kubestone/controllers/fio"
	"github.com/xridge/kubestone/controllers/iperf3"
	"github.com/xridge/kubestone/controllers/sysbench"
	"github.com/xridge/kubestone/pkg/k8s"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	rootLog  = zap.RawLoggerTo(os.Stderr, true)
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = k8sscheme.AddToScheme(scheme)

	_ = perfv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func newEventRecorder(clientSet *kubernetes.Clientset) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(rootLog.Sugar().Infof)
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{Interface: clientSet.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(k8sscheme.Scheme,
		corev1.EventSource{Component: "kubestone"})
	return recorder
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zapr.NewLogger(rootLog))

	restClientConfig := ctrl.GetConfigOrDie()
	mgr, err := ctrl.NewManager(restClientConfig, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
	})
	if err != nil {
		setupLog.Error(err, "Unable to start manager")
		os.Exit(1)
	}

	clientSet := kubernetes.NewForConfigOrDie(restClientConfig)
	k8sAccess := k8s.Access{
		Client:        mgr.GetClient(),
		Clientset:     clientSet,
		Scheme:        mgr.GetScheme(),
		EventRecorder: newEventRecorder(clientSet),
	}

	if err = (&iperf3.Reconciler{
		K8S: k8sAccess,
		Log: ctrl.Log.WithName("controllers").WithName("Iperf3"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Iperf3")
		os.Exit(1)
	}
	if err = (&fio.Reconciler{
		K8S: k8sAccess,
		Log: ctrl.Log.WithName("controllers").WithName("Fio"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Fio")
		os.Exit(1)
	}
	if err = (&sysbench.Reconciler{
		K8S: k8sAccess,
		Log: ctrl.Log.WithName("controllers").WithName("Sysbench"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Sysbench")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
