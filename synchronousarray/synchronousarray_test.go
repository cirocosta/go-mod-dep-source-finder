package synchronous_array_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cirocosta/go-mod-license-finder/result"
	. "github.com/cirocosta/go-mod-license-finder/synchronousarray"
)

var _ = Describe("SynchronousArray", func() {
	Context("When adding items from many routines", func() {

		var (
			array        SynchronousArray
			resultsToAdd []Result
		)

		BeforeEach(func() {
			array = SynchronousArray{}
			resultsToAdd = []Result{
				{
					Original: "1",
				},
				{
					Original: "2",
				},
				{
					Original: "3",
				},
			}
		})

		JustBeforeEach(func() {
			for _, resultToAdd := range resultsToAdd {
				go func(result Result) {
					array.Add(result)
				}(resultToAdd)
			}
		})

		It("Adds all items", func() {
			Eventually(func() []Result {
				return array.Results
			}).Should(HaveLen(3))

			Expect(array.Results).Should(ConsistOf(resultsToAdd))
		})
	})
})
