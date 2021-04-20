package types_test

import (
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/ginkgo/types"
)

var _ = Describe("Types", func() {
	Describe("Report", func() {
		Describe("Add", func() {
			It("concatenates spec reports, combines success, and computes a new RunTime", func() {
				t := time.Now()
				reportA := types.Report{
					SuitePath:                 "foo",
					SuiteSucceeded:            true,
					StartTime:                 t.Add(-time.Minute),
					EndTime:                   t.Add(2 * time.Minute),
					SpecialSuiteFailureReason: "",
					SpecReports: types.SpecReports{
						types.SpecReport{NumAttempts: 3},
						types.SpecReport{NumAttempts: 4},
					},
				}

				reportB := types.Report{
					SuitePath:                 "bar",
					SuiteSucceeded:            false,
					StartTime:                 t.Add(-2 * time.Minute),
					EndTime:                   t.Add(time.Minute),
					SpecialSuiteFailureReason: "blame bob",
					SpecReports: types.SpecReports{
						types.SpecReport{NumAttempts: 5},
						types.SpecReport{NumAttempts: 6},
					},
				}

				composite := reportA.Add(reportB)
				Ω(composite).Should(Equal(types.Report{
					SuitePath:                 "foo",
					SuiteSucceeded:            false,
					StartTime:                 t.Add(-2 * time.Minute),
					EndTime:                   t.Add(2 * time.Minute),
					RunTime:                   4 * time.Minute,
					SpecialSuiteFailureReason: "blame bob",
					SpecReports: types.SpecReports{
						types.SpecReport{NumAttempts: 3},
						types.SpecReport{NumAttempts: 4},
						types.SpecReport{NumAttempts: 5},
						types.SpecReport{NumAttempts: 6},
					},
				}))

			})
		})
	})

	Describe("NodeType", func() {
		Describe("Is", func() {
			It("returns true when the NodeType is in the passed-in list", func() {
				Ω(types.NodeTypeContainer.Is(types.NodeTypeIt, types.NodeTypeContainer)).Should(BeTrue())
			})

			It("returns false when the NodeType is not in the passed-in list", func() {
				Ω(types.NodeTypeContainer.Is(types.NodeTypeIt, types.NodeTypeBeforeEach)).Should(BeFalse())
			})
		})

		DescribeTable("Representation and Encoding", func(nodeType types.NodeType, expectedString string) {
			Ω(nodeType.String()).Should(Equal(expectedString))

			marshalled, err := json.Marshal(nodeType)
			Ω(err).ShouldNot(HaveOccurred())
			var unmarshalled types.NodeType
			json.Unmarshal(marshalled, &unmarshalled)
			Ω(unmarshalled).Should(Equal(nodeType))
		},
			Entry("Container", types.NodeTypeContainer, "Container"),
			Entry("It", types.NodeTypeIt, "It"),
			Entry("BeforeEach", types.NodeTypeBeforeEach, "BeforeEach"),
			Entry("JustBeforeEach", types.NodeTypeJustBeforeEach, "JustBeforeEach"),
			Entry("AfterEach", types.NodeTypeAfterEach, "AfterEach"),
			Entry("JustAfterEach", types.NodeTypeJustAfterEach, "JustAfterEach"),
			Entry("BeforeSuite", types.NodeTypeBeforeSuite, "BeforeSuite"),
			Entry("SynchronizedBeforeSuite", types.NodeTypeSynchronizedBeforeSuite, "SynchronizedBeforeSuite"),
			Entry("AfterSuite", types.NodeTypeAfterSuite, "AfterSuite"),
			Entry("SynchronizedAfterSuite", types.NodeTypeSynchronizedAfterSuite, "SynchronizedAfterSuite"),
			Entry("Invalid", types.NodeTypeInvalid, "INVALID NODE TYPE"),
		)
	})

	Describe("FailureNodeContext", func() {
		DescribeTable("Representation and Encoding", func(context types.FailureNodeContext) {
			marshalled, err := json.Marshal(context)
			Ω(err).ShouldNot(HaveOccurred())
			var unmarshalled types.FailureNodeContext
			json.Unmarshal(marshalled, &unmarshalled)
			Ω(unmarshalled).Should(Equal(context))
		},
			Entry("LeafNode", types.FailureNodeIsLeafNode),
			Entry("TopLevel", types.FailureNodeAtTopLevel),
			Entry("InContainer", types.FailureNodeInContainer),
			Entry("Invalid", types.FailureNodeContextInvalid),
		)
	})

	Describe("SpecState", func() {
		DescribeTable("Representation and Encoding", func(specState types.SpecState, expectedString string) {
			Ω(specState.String()).Should(Equal(expectedString))

			marshalled, err := json.Marshal(specState)
			Ω(err).ShouldNot(HaveOccurred())
			var unmarshalled types.SpecState
			json.Unmarshal(marshalled, &unmarshalled)
			Ω(unmarshalled).Should(Equal(specState))
		},
			Entry("Pending", types.SpecStatePending, "pending"),
			Entry("Skipped", types.SpecStateSkipped, "skipped"),
			Entry("Passed", types.SpecStatePassed, "passed"),
			Entry("Failed", types.SpecStateFailed, "failed"),
			Entry("Panicked", types.SpecStatePanicked, "panicked"),
			Entry("Interrupted", types.SpecStateInterrupted, "interrupted"),
			Entry("Invalid", types.SpecStateInvalid, "INVALID SPEC STATE"),
		)
	})

	Describe("SpecReport Helper Functions", func() {
		Describe("CombinedOutput", func() {
			Context("with no GinkgoWriter or StdOutErr output", func() {
				It("comes back empty", func() {
					Ω(types.SpecReport{}.CombinedOutput()).Should(Equal(""))
				})
			})

			Context("wtih only StdOutErr output", func() {
				It("returns that output", func() {
					Ω(types.SpecReport{
						CapturedStdOutErr: "hello",
					}.CombinedOutput()).Should(Equal("hello"))
				})
			})

			Context("wtih only GinkgoWriter output", func() {
				It("returns that output", func() {
					Ω(types.SpecReport{
						CapturedGinkgoWriterOutput: "hello",
					}.CombinedOutput()).Should(Equal("hello"))
				})
			})

			Context("wtih both", func() {
				It("returns both concatenated", func() {
					Ω(types.SpecReport{
						CapturedGinkgoWriterOutput: "gw",
						CapturedStdOutErr:          "std",
					}.CombinedOutput()).Should(Equal("std\ngw"))
				})
			})
		})

		It("can report on whether state is a failed state", func() {
			Ω(types.SpecReport{State: types.SpecStatePending}.Failed()).Should(BeFalse())
			Ω(types.SpecReport{State: types.SpecStateSkipped}.Failed()).Should(BeFalse())
			Ω(types.SpecReport{State: types.SpecStatePassed}.Failed()).Should(BeFalse())
			Ω(types.SpecReport{State: types.SpecStateFailed}.Failed()).Should(BeTrue())
			Ω(types.SpecReport{State: types.SpecStatePanicked}.Failed()).Should(BeTrue())
			Ω(types.SpecReport{State: types.SpecStateInterrupted}.Failed()).Should(BeTrue())
		})

		It("can return a concatenated set of texts", func() {
			Ω(CurrentSpecReport().FullText()).Should(Equal("Types SpecReport Helper Functions can return a concatenated set of texts"))
		})

		It("can return the name of the file it's spec is in", func() {
			cl := types.NewCodeLocation(0)
			Ω(CurrentSpecReport().FileName()).Should(Equal(cl.FileName))
		})

		It("can return the linenumber of the file it's spec is in", func() {
			cl := types.NewCodeLocation(0)
			Ω(CurrentSpecReport().LineNumber()).Should(Equal(cl.LineNumber - 1))
		})

		It("can return it's failure's message", func() {
			report := types.SpecReport{
				Failure: types.Failure{Message: "why this failed"},
			}
			Ω(report.FailureMessage()).Should(Equal("why this failed"))
		})

		It("can return it's failure's code location", func() {
			cl := types.NewCodeLocation(0)
			report := types.SpecReport{
				Failure: types.Failure{Location: cl},
			}
			Ω(report.FailureLocation()).Should(Equal(cl))
		})
	})

	Describe("SpecReports", func() {
		Describe("Encoding to JSON", func() {
			var report types.SpecReport
			BeforeEach(func() {
				report = types.SpecReport{
					ContainerHierarchyTexts: []string{"A", "B"},
					ContainerHierarchyLocations: []types.CodeLocation{
						types.NewCodeLocation(0),
						types.NewCodeLocationWithStackTrace(0),
						types.NewCustomCodeLocation("welp"),
					},
					LeafNodeType:               types.NodeTypeIt,
					LeafNodeLocation:           types.NewCodeLocation(0),
					LeafNodeText:               "C",
					State:                      types.SpecStateFailed,
					StartTime:                  time.Date(2012, 06, 19, 05, 32, 12, 0, time.UTC),
					EndTime:                    time.Date(2012, 06, 19, 05, 33, 12, 0, time.UTC),
					RunTime:                    time.Minute,
					GinkgoParallelNode:         2,
					NumAttempts:                3,
					CapturedGinkgoWriterOutput: "gw",
					CapturedStdOutErr:          "std",
					Failure: types.Failure{
						Message:                   "boom",
						Location:                  types.NewCodeLocation(1),
						ForwardedPanic:            "bam",
						FailureNodeContext:        types.FailureNodeInContainer,
						FailureNodeType:           types.NodeTypeBeforeEach,
						FailureNodeLocation:       types.NewCodeLocation(0),
						FailureNodeContainerIndex: 1,
					},
				}
			})

			Context("with a failure", func() {
				It("round-trips correctly", func() {
					marshalled, err := json.Marshal(report)
					Ω(err).ShouldNot(HaveOccurred())
					unmarshalled := types.SpecReport{}
					err = json.Unmarshal(marshalled, &unmarshalled)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(unmarshalled).Should(Equal(report))
				})
			})

			Context("without a failure", func() {
				BeforeEach(func() {
					report.Failure = types.Failure{}
				})
				It("round-trips correclty and doesn't include the Failure struct", func() {
					marshalled, err := json.Marshal(report)
					Ω(string(marshalled)).ShouldNot(ContainSubstring("Failure"))
					Ω(err).ShouldNot(HaveOccurred())
					unmarshalled := types.SpecReport{}
					err = json.Unmarshal(marshalled, &unmarshalled)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(unmarshalled).Should(Equal(report))
				})
			})
		})

		Describe("WithLeafNodeType", func() {
			It("returns reports with the matching LeafNodeTypes", func() {
				reports := types.SpecReports{
					{LeafNodeType: types.NodeTypeIt, NumAttempts: 2},
					{LeafNodeType: types.NodeTypeIt, NumAttempts: 3},
					{LeafNodeType: types.NodeTypeBeforeSuite, NumAttempts: 4},
					{LeafNodeType: types.NodeTypeAfterSuite, NumAttempts: 5},
					{LeafNodeType: types.NodeTypeSynchronizedAfterSuite, NumAttempts: 6},
				}

				Ω(reports.WithLeafNodeType(types.NodeTypeIt, types.NodeTypeAfterSuite)).Should(Equal(types.SpecReports{
					{LeafNodeType: types.NodeTypeIt, NumAttempts: 2},
					{LeafNodeType: types.NodeTypeIt, NumAttempts: 3},
					{LeafNodeType: types.NodeTypeAfterSuite, NumAttempts: 5},
				}))
			})
		})

		Describe("WithState", func() {
			It("returns reports with the matching SpecStates", func() {
				reports := types.SpecReports{
					{State: types.SpecStatePassed, NumAttempts: 2},
					{State: types.SpecStatePassed, NumAttempts: 3},
					{State: types.SpecStateFailed, NumAttempts: 4},
					{State: types.SpecStatePending, NumAttempts: 5},
					{State: types.SpecStateSkipped, NumAttempts: 6},
				}

				Ω(reports.WithState(types.SpecStatePassed, types.SpecStatePending)).Should(Equal(types.SpecReports{
					{State: types.SpecStatePassed, NumAttempts: 2},
					{State: types.SpecStatePassed, NumAttempts: 3},
					{State: types.SpecStatePending, NumAttempts: 5},
				}))
			})
		})

		Describe("CountWithState", func() {
			It("returns the number with the matching SpecStates", func() {
				reports := types.SpecReports{
					{State: types.SpecStatePassed, NumAttempts: 2},
					{State: types.SpecStatePassed, NumAttempts: 3},
					{State: types.SpecStateFailed, NumAttempts: 4},
					{State: types.SpecStatePending, NumAttempts: 5},
					{State: types.SpecStateSkipped, NumAttempts: 6},
				}

				Ω(reports.CountWithState(types.SpecStatePassed, types.SpecStatePending)).Should(Equal(3))
			})
		})

		Describe("CountOfFlakedSpecs", func() {
			It("returns the number of passing specs with NumAttempts > 1", func() {
				reports := types.SpecReports{
					{State: types.SpecStatePassed, NumAttempts: 2},
					{State: types.SpecStatePassed, NumAttempts: 2},
					{State: types.SpecStatePassed, NumAttempts: 1},
					{State: types.SpecStatePassed, NumAttempts: 1},
					{State: types.SpecStateFailed, NumAttempts: 2},
				}

				Ω(reports.CountOfFlakedSpecs()).Should(Equal(2))
			})
		})
	})
})
