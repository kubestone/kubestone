package e2e

import (
	"bytes"
	"context"
	"log"
	"os/exec"

	"github.com/firepear/qsplit"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xridge/kubestone/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

const (
	iperf3SampleCR = "../../config/samples/perf_v1alpha1_iperf3.yaml"
	e2eNamespace   = "kubestone-e2e"
)

var restClientConfig = ctrl.GetConfigOrDie()
var client ctrlclient.Client
var ctx = context.Background()
var scheme = runtime.NewScheme()

func init() {
	_ = k8sscheme.AddToScheme(scheme)
	err := perfv1alpha1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	client, err = ctrlclient.New(restClientConfig, ctrlclient.Options{Scheme: scheme})
	if err != nil {
		panic(err)
	}
}

func run(command string) (stdout, stderr string, err error) {
	commandArray := qsplit.ToStrings([]byte(command))
	cmd := exec.Command(commandArray[0], commandArray[1:]...)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdOut, &stdErr
	err = cmd.Run()
	if err != nil {
		log.Printf("Error during execution of `%v`\nerr: %v\nstdout: %v\nstderr: %v\n",
			command, err, stdOut.String(), stdErr.String())
		return "", "", err
	}

	return stdOut.String(), stdErr.String(), nil
}

var _ = Describe("end to end test", func() {
	Context("preparing namespace", func() {
		_, _, err := run("kubectl create namespace " + e2eNamespace)
		It("should succeed", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("for iperf3", func() {
		Context("creation from samples", func() {
			_, _, err := run("kubectl create -n " + e2eNamespace + " -f " + iperf3SampleCR)
			It("should create iperf3-sample cr", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("created job", func() {
			It("Should finish in a pre-defined time", func() {
				timeout := 60
				cr := &v1alpha1.Iperf3{}
				namespacedName := types.NamespacedName{
					Namespace: e2eNamespace,
					Name:      "iperf3-sample",
				}
				Eventually(func() bool {
					if err := client.Get(ctx, namespacedName, cr); err != nil {
						panic(err)
					}
					return (cr.Status.Running == false) && (cr.Status.Completed)
				}, timeout).Should(BeTrue())
			})
			It("Should leave one successful pod", func() {
				pod := &corev1.Pod{}
				namespacedName := types.NamespacedName{
					Namespace: e2eNamespace,
					Name:      "iperf3-sample-client",
				}
				Expect(client.Get(ctx, namespacedName, pod)).To(Succeed())
			})
			It("Should not leave deployment", func() {
				deployment := &appsv1.Deployment{}
				namespacedName := types.NamespacedName{
					Namespace: e2eNamespace,
					Name:      "iperf3-sample",
				}
				Expect(client.Get(ctx, namespacedName, deployment)).NotTo(Succeed())
			})
			It("Should not leave service", func() {
				service := &corev1.Service{}
				namespacedName := types.NamespacedName{
					Namespace: e2eNamespace,
					Name:      "iperf3-sample",
				}
				Expect(client.Get(ctx, namespacedName, service)).NotTo(Succeed())
			})
		})
	})
	Describe("for fio", func() {
		// FIXME: Add fio e2e
	})
})
