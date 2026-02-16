package orderhdl_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOrderhdl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Orderhdl Suite")
}
