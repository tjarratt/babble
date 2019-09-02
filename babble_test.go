package babble_test

import (
	. "github.com/tjarratt/babble"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("babble", func() {
	var babbler Babbler
	BeforeEach(func() {
		babbler = Babbler{
			Count: 1,
			Words: []string{"hello"},
			Separator: "☃",
		}
	})

	It("returns a random word", func() {
		Expect(babbler.Babble()).To(Equal("hello"))
	})

	Describe("with multiple words", func() {
		It("concatenates strings", func() {
			babbler.Count = 2
			Expect(babbler.Babble()).To(Equal("hello☃hello"))
		})
	})
})
