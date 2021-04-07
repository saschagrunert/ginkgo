package internal_integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/internal/test_helpers"
	. "github.com/onsi/gomega"
)

var _ = Describe("when config.FailFast is enabled", func() {
	BeforeEach(func() {
		conf.FailFast = true

		RunFixture("fail fast", func() {
			Describe("a container", func() {
				BeforeEach(rt.T("bef"))
				It("A", rt.T("A"))
				It("B", rt.T("B", func() { F() }))
				It("C", rt.T("C", func() { F() }))
				It("D", rt.T("D"))
				AfterEach(rt.T("aft"))
			})
			AfterSuite(rt.T("after-suite"))
		})
	})

	It("does not run any tests after the failure occurs, but does run the failed tests's after each and the after suite", func() {
		Ω(rt).Should(HaveTracked(
			"bef", "A", "aft",
			"bef", "B", "aft",
			"after-suite",
		))
	})

	It("reports that the tests were skipped", func() {
		Ω(reporter.Did.Find("A")).Should(HavePassed())
		Ω(reporter.Did.Find("B")).Should(HaveFailed())
		Ω(reporter.Did.Find("C")).Should(HaveBeenSkipped())
		Ω(reporter.Did.Find("D")).Should(HaveBeenSkipped())
	})

	It("reports the correct statistics", func() {
		Ω(reporter.End).Should(BeASuiteSummary(NSpecs(4), NPassed(1), NFailed(1), NSkipped(2)))
	})
})
