package internal_integration_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/internal/test_helpers"
	"github.com/onsi/ginkgo/types"
	. "github.com/onsi/gomega"
)

var _ = Describe("When a test suite is interrupted", func() {
	Describe("when it is interrupted in a BeforeSuite", func() {
		BeforeEach(func() {
			success, _ := RunFixture("interrupted test", func() {
				BeforeSuite(rt.T("before-suite", func() {
					interruptHandler.Interrupt()
					time.Sleep(time.Hour)
				}))
				AfterSuite(rt.T("after-suite"))
				It("A", rt.T("A"))
				It("B", rt.T("B"))
			})
			Ω(success).Should(Equal(false))
		})

		It("runs the AfterSuite and skips all the tests", func() {
			Ω(rt).Should(HaveTracked("before-suite", "after-suite"))
			Ω(reporter.Did.FindByLeafNodeType(types.NodeTypeIt)).Should(BeZero())
		})

		It("reports the correct failure", func() {
			summary := reporter.Did.FindByLeafNodeType(types.NodeTypeBeforeSuite)
			Ω(summary.State).Should(Equal(types.SpecStateInterrupted))
			Ω(summary.Failure.Message).Should(ContainSubstring("Interrupted by User"))
		})

		It("reports the correct statistics", func() {
			Ω(reporter.End).Should(BeASuiteSummary(false, NSpecs(2), NWillRun(2), NPassed(0), NSkipped(2), NFailed(0)))
		})
	})

	Describe("when it is interrupted in a test", func() {
		BeforeEach(func() {
			conf.FlakeAttempts = 3
			success, _ := RunFixture("interrupted test", func() {
				BeforeSuite(rt.T("before-suite"))
				AfterSuite(rt.T("after-suite"))
				BeforeEach(rt.T("bef.1"))
				AfterEach(rt.T("aft.1"))
				Describe("container", func() {
					BeforeEach(rt.T("bef.2"))
					AfterEach(rt.T("aft.2"))
					It("runs", rt.T("runs"))
					Describe("nested-container", func() {
						BeforeEach(rt.T("bef.3-interrupt!", func() {
							interruptHandler.Interrupt()
							time.Sleep(time.Hour)
						}))
						AfterEach(rt.T("aft.3a"))
						AfterEach(rt.T("aft.3b", func() {
							interruptHandler.Interrupt()
							time.Sleep(time.Hour)
						}))
						Describe("deeply-nested-container", func() {
							BeforeEach(rt.T("bef.4"))
							AfterEach(rt.T("aft.4"))
							It("the interrupted test", rt.T("the interrupted test"))
							It("skipped.1", rt.T("skipped.1"))
						})
					})
					It("skipped.2", rt.T("skipped.2"))
				})
			})
			Ω(success).Should(Equal(false))
		})

		It("unwinds the after eaches at the appropriate nesting level, allowing additional interrupts of after eaches as it goes", func() {
			Ω(rt).Should(HaveTracked("before-suite",
				"bef.1", "bef.2", "runs", "aft.2", "aft.1",
				"bef.1", "bef.2", "bef.3-interrupt!", "aft.3a", "aft.3b", "aft.2", "aft.1",
				"after-suite"))
		})

		It("skips subsequent tests", func() {
			Ω(reporter.Did.WithState(types.SpecStatePassed).Names()).Should(ConsistOf("runs"))
			Ω(reporter.Did.WithState(types.SpecStateInterrupted).Names()).Should(ConsistOf("the interrupted test"))
			Ω(reporter.Did.WithState(types.SpecStateSkipped).Names()).Should(ConsistOf("skipped.1", "skipped.2"))
		})

		It("reports the interrupted test as interrupted and emits a stack trace", func() {
			message := reporter.Did.Find("the interrupted test").Failure.Message
			Ω(message).Should(ContainSubstring("Interrupted by User"))
			Ω(message).Should(ContainSubstring("Here's a stack trace of all running goroutines:"))
			Ω(message).Should(ContainSubstring("internal.interruptMessageWithStackTraces"))
		})

		It("reports the correct statistics", func() {
			Ω(reporter.End).Should(BeASuiteSummary(false, NSpecs(4), NWillRun(4), NPassed(1), NSkipped(2), NFailed(1)))
		})
	})
})
