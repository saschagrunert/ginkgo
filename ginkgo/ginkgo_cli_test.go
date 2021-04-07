package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Command Documentation Links", func() {
	It("all commands with doc links point to valid documentation", func() {
		commands := GenerateCommands()
		for _, command := range commands {
			if command.DocLink != "" {
				Ω(command.DocLink).Should(BeElementOf(DOC_ANCHORS))
			}
		}
	})
})
