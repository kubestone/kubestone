package fio

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFioController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fio Controller Suite")
}
