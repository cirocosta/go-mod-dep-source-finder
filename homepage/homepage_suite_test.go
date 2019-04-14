package homepage_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHomepage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Homepage Suite")
}
