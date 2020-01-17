package ycsbbench

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestYcsbBenchController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "YcsbBench Controller Suite")
}
