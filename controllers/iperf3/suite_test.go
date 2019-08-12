package iperf3

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIperf3Controller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Iperf3 Controller Suite")
}
