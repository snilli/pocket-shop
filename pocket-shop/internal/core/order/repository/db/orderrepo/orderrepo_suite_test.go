package orderrepo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOrderrepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Orderrepo Integration Suite")
}
