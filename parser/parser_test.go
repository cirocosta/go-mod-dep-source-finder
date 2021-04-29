package parser_test

import (
	"github.com/cirocosta/go-mod-dep-source-finder/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("StripVersionFromDependency", func() {

	type testCase struct {
		input       string
		expected    string
		shouldError bool
	}

	DescribeTable("varying possible dependencies",
		func(tc testCase) {
			actual, err := parser.StripVersionFromDependency(tc.input)

			if tc.shouldError {
				Expect(err).To(HaveOccurred())
				return
			}

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(tc.expected))
		},
		Entry("empty", testCase{
			input:       "",
			shouldError: true,
		}),
		Entry("without version at the end", testCase{
			input:       "github.com/something/else",
			shouldError: true,
		}),
		Entry("with version at the end", testCase{
			input:    "github.com/something/else/v2",
			expected: "github.com/something/else",
		}),
	)

})

var _ = Describe("ParseLine", func() {

	type testCase struct {
		line        string
		parsedLine  parser.Line
		shouldError bool
	}

	DescribeTable("varying possible lines",
		func(tc testCase) {
			line, err := parser.ParseLine(tc.line)

			if tc.shouldError {
				Expect(err).To(HaveOccurred())
				return
			}

			Expect(err).ToNot(HaveOccurred())
			Expect(line).To(Equal(tc.parsedLine))

		},
		Entry("empty", testCase{
			line:        "",
			shouldError: true,
		}),
		Entry("without enough fields", testCase{
			line:        "aaaaa",
			shouldError: true,
		}),
		Entry("without a semver after first field", testCase{
			line:        "aaa bbb",
			shouldError: true,
		}),
		Entry("having a proper semver after the dependency", testCase{
			line:       "aaa v1.2.3",
			parsedLine: parser.Line{"aaa", "v1.2.3", ""},
		}),
		Entry("having a replace directive", testCase{
			line:       "aaa v1.2.3 => bbb v1.22.33",
			parsedLine: parser.Line{"bbb", "v1.22.33", ""},
		}),
		Entry("having a replace to existing relative directory", testCase{
			line:       "github.com/xoebus/go-tracker v0.0.0-00010101000000-000000000000 => ./testdata/",
			parsedLine: parser.Line{"", "", "./testdata/"},
		}),
		Entry("having a replace to inexistent relative directory", testCase{
			line:        "github.com/xoebus/go-tracker v0.0.0-00010101000000-000000000000 => ./inexistent/",
			shouldError: true,
		}),
		Entry("having leading spaces", testCase{
			line:       "   aaa v1.2.3",
			parsedLine: parser.Line{"aaa", "v1.2.3", ""},
		}),
		Entry("having version suffix", testCase{
			line:       "   aaa/v3 v3.2.3",
			parsedLine: parser.Line{"aaa", "v3.2.3", ""},
		}),
		Entry("being an incompatible version", testCase{
			line:       "code.cloudfoundry.org/lager v2.0.0+incompatible",
			parsedLine: parser.Line{"code.cloudfoundry.org/lager", "v2.0.0", ""},
		}),
		Entry("being a directory within a github repo", testCase{
			line:       "github.com/ugorji/go/codec v0.0.0-20181209151446-772ced7fd4c2 // indirect",
			parsedLine: parser.Line{"github.com/ugorji/go", "772ced7fd4c2", ""},
		}),
		Entry("having trailing fields", testCase{
			line:       "   aaa v1.2.3 // indirect",
			parsedLine: parser.Line{"aaa", "v1.2.3", ""},
		}),
		Entry("having a pseudo-version", testCase{
			line:       "   aaa v0.0.0-20180518195852-02e53af36e6c // indirect",
			parsedLine: parser.Line{"aaa", "02e53af36e6c", ""},
		}),
	)
})

var _ = Describe("ParseVersion", func() {

	type testCase struct {
		version     string
		reference   string
		shouldError bool
	}

	DescribeTable("varying the possible versions",
		func(tc testCase) {
			reference, err := parser.ParseVersion(tc.version)

			if tc.shouldError {
				Expect(err).To(HaveOccurred())
				return
			}

			Expect(err).ToNot(HaveOccurred())
			Expect(reference).To(Equal(tc.reference))
		},
		Entry("empty", testCase{
			version:     "",
			shouldError: true,
		}),
		Entry("not being a semver tag", testCase{
			version:     "udshudhk",
			shouldError: true,
		}),
		Entry("being a plain-version semver", testCase{
			version:   "v1.2.3",
			reference: "v1.2.3",
		}),
		Entry("having semver followed by date and sha", testCase{
			version:   "v0.0.0-20180518195852-02e53af36e6c",
			reference: "02e53af36e6c",
		}),
		Entry("having semver followed by date and sha and incompatible", testCase{
			version:   "v2.0.0-alpha.0.0.20171101191150-72e1c2a1ef30+incompatible",
			reference: "72e1c2a1ef30",
		}),
		Entry("having server with +incompatible", testCase{
			version:   "v2.0.0+incompatible",
			reference: "v2.0.0",
		}),
	)
})
