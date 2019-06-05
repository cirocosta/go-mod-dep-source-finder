package synchronous_array_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSynchronousArray(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SynchronousArray Suite")
}
