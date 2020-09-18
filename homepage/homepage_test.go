package homepage_test

import (
	"github.com/cirocosta/go-mod-dep-source-finder/homepage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Homepage", func() {

	type testCase struct {
		dependency  string
		version     string
		expected    string
		unknownHost bool
		shouldError bool
	}

	DescribeTable("varying possible dependencies",
		func(tc testCase) {
			actual, unknownHost, err := homepage.Find(tc.dependency, tc.version)
			if tc.shouldError {
				Expect(err).To(HaveOccurred())
				return
			}

			Expect(err).ToNot(HaveOccurred())
			Expect(unknownHost).To(Equal(tc.unknownHost))
			Expect(actual).To(Equal(tc.expected))

		},
		Entry("unknown host", testCase{
			dependency:  "https://something.else/pkg",
			version:     "v2.0.0",
			unknownHost: true,
		}),
		Entry("gitlab tag", testCase{
			dependency: "https://gitlab.com/cznic/golex",
			version:    "v1.0.0",
			expected:   "https://gitlab.com/cznic/golex/tree/v1.0.0",
		}),
		Entry("github tag", testCase{
			dependency: "https://github.com/cloudfoundry/urljoiner",
			version:    "v1.2.3",
			expected:   "https://github.com/cloudfoundry/urljoiner/tree/v1.2.3",
		}),
		Entry("github sha", testCase{
			dependency: "https://github.com/cloudfoundry/urljoiner",
			version:    "fa338ed9e9ec",
			expected:   "https://github.com/cloudfoundry/urljoiner/tree/fa338ed9e9ec",
		}),
		Entry("bitbucket sha", testCase{
			dependency: "https://bitbucket.org/foo/bar",
			version:    "fa338ed9e9ec",
			expected:   "https://bitbucket.org/foo/bar/src/fa338ed9e9ec",
		}),
		Entry("bitbucket tag", testCase{
			dependency: "https://bitbucket.org/foo/bar",
			version:    "v2.0.0",
			expected:   "https://bitbucket.org/foo/bar/src/v2.0.0",
		}),
		Entry("gopkg.in without user", testCase{
			dependency: "https://gopkg.in/yaml.v2",
			version:    "v2.2.2",
			expected:   "https://github.com/go-yaml/yaml/tree/v2.2.2",
		}),
		Entry("gopkg.in with user", testCase{
			dependency: "https://gopkg.in/ory-am/dockertest.v2",
			version:    "v2.2.3",
			expected:   "https://github.com/ory-am/dockertest/tree/v2.2.3",
		}),
		Entry("go.googlesource.com tag", testCase{
			dependency: "https://go.googlesource.com/text",
			version:    "v0.3.0",
			expected:   "https://go.googlesource.com/text/+/v0.3.0/",
		}),
		Entry("code.googlesource.com tag", testCase{
			dependency: "https://code.googlesource.com/google-api-go-client",
			version:    "v0.3.0",
			expected:   "https://code.googlesource.com/google-api-go-client/+/v0.3.0/",
		}),
	)
})
