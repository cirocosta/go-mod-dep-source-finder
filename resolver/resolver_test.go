package resolver_test

import (
	"bytes"
	"context"

	"github.com/cirocosta/go-mod-license-finder/parser"
	"github.com/cirocosta/go-mod-license-finder/resolver"
	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FindGoImport", func() {

	var (
		content    string
		importLine string
		found      bool
		err        error
	)

	JustBeforeEach(func() {
		importLine, found, err = resolver.FindGoImport(bytes.NewBufferString(content))
	})

	Context("providing invalid html", func() {

		BeforeEach(func() {
			content = `<html><body thiiiis>>`
		})

		It("doesn't error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("doesn't find anything", func() {
			Expect(found).ToNot(BeTrue())
		})

	})

	Context("having a proper html", func() {

		Context("not having any go-import in the html", func() {
			BeforeEach(func() {
				content = `<html><body></body></html>`
			})

			It("doesn't find", func() {
				Expect(found).ToNot(BeTrue())
			})
		})

		Context("having a go-import without content", func() {
			BeforeEach(func() {
				content = `<html><meta name="go-import" content=""></html>`
			})

			It("finds returning empty", func() {
				Expect(found).To(BeTrue())
				Expect(importLine).To(BeEmpty())
			})
		})
	})
})

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
