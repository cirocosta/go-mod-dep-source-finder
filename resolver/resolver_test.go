package resolver_test

import (
	"bytes"
	"context"

	"github.com/cirocosta/go-mod-license-finder/resolver"
	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseGoImportContent", func() {

	type testCase struct {
		content     string
		goImport    resolver.GoImport
		shouldError bool
	}

	DescribeTable("varying possible content",
		func(tc testCase) {
			res, err := resolver.ParseGoImport(tc.content)
			if tc.shouldError {
				Expect(err).To(HaveOccurred())
				return
			}

			Expect(res).To(Equal(tc.goImport))
		},
		Entry("empty", testCase{
			content:     "",
			shouldError: true,
		}),
		Entry("without enough fields", testCase{
			content:     "github.com/cirocosta/l4 git",
			shouldError: true,
		}),
		Entry("with enough fields", testCase{
			content: "github.com/cirocosta/l4 git https://github.com/cirocosta/l4.git",
			goImport: resolver.GoImport{
				ImportPrefix: "github.com/cirocosta/l4",
				VCS:          "git",
				RepoRoot:     "https://github.com/cirocosta/l4.git",
			},
		}),
	)

})

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

var _ = Describe("Resolve", func() {

	var (
		server     *ghttp.Server
		dependency string
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
			dependency = server.URL()

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/"),
			))
		})

		It("issues request to the domain", func() {
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})
})
