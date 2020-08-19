package kafkabench

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKafkaBenchController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KafkaBench Controller Suite")
}
