package babble_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBabble(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Babble Suite")
}
