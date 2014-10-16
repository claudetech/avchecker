package avchecker

import (
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAvailabilityChecker(t *testing.T) {
	RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "AvailabilityChecker Suite")
}
