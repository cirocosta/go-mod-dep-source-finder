package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/cirocosta/go-mod-license-finder/parser"
)

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
			line:       "aaa bbb",
			parsedLine: parser.Line{"aaa", "bbb"},
		}),
		Entry("having leading white spaces", testCase{
			line: "   	aaa bbb",
			parsedLine: parser.Line{"aaa", "bbb"},
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
