package resolver_test

import (
	"context"

	"github.com/cirocosta/go-mod-license-finder/parser"
	"github.com/cirocosta/go-mod-license-finder/resolver"
	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resolver", func() {

	var (
		server     *ghttp.Server
		dependency parser.Line
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
	})

	AfterEach(func() {
		server.Close()
	})

	JustBeforeEach(func() {
		_, _ = resolver.Resolve(context.TODO(), dependency)
	})

	Context("having a proper dependency as input", func() {

		BeforeEach(func() {
			dependency = parser.Line{
				Dependency: server.URL(),
				Reference:  "v1.2.3",
			}

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/"),
			))
		})

		It("issues request to the domain", func() {
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})
})
